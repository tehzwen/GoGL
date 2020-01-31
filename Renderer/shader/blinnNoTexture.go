package shader

type BlinnNoTexture struct {
	fragShader string
	vertShader string
	geoShader  string
}

func (s BlinnNoTexture) GetFragShader() string {
	return s.fragShader
}

func (s BlinnNoTexture) GetVertShader() string {
	return s.vertShader
}

func (s BlinnNoTexture) GetGeometryShader() string {
	return s.geoShader
}

func (s *BlinnNoTexture) Setup() {
	s.vertShader = `
	#version 410
	//needed to add layout location for mac to work properly
	layout (location = 0) in vec3 aPosition;
	layout (location = 1) in vec3 aNormal;

	out vec3 oNormal;
	out vec3 normalInterp;
	out vec3 oFragPosition;
	out vec3 oCamPosition;
	
	uniform vec3 cameraPosition;
	uniform mat4 uProjectionMatrix;
	uniform mat4 uViewMatrix;
	uniform mat4 uModelMatrix;

	void main() {
		mat4 normalMatrix = transpose(inverse(uModelMatrix));
		oNormal = normalize((uModelMatrix * vec4(aNormal, 1.0)).xyz);
		normalInterp = vec3(normalMatrix * vec4(aNormal, 0.0));
		oFragPosition = (uModelMatrix * vec4(aPosition, 1.0)).xyz;
		oCamPosition =  (uViewMatrix * vec4(cameraPosition, 1.0)).xyz;
		gl_Position = uProjectionMatrix * uViewMatrix * uModelMatrix * vec4(aPosition, 1.0); 
	}
` + "\x00"

	s.geoShader = ""

	s.fragShader = `
	#version 410
	precision highp float;
	#define MAX_LIGHTS 128

	struct PointLight {
		vec3 position;
		float strength;
		float constant;
		float linear;
		float quadratic; 
		vec3 color;
	};

	in vec3 oFragPosition;
	in vec3 normalInterp;
	in vec3 oNormal;
	in vec3 oCamPosition;
	
	uniform vec3 diffuseVal;
	uniform vec3 ambientVal;
	uniform vec3 specularVal;
	uniform float nVal;
	uniform float Alpha;
	uniform int numLights;
	uniform samplerCube depthMap;
	uniform PointLight pointLights[MAX_LIGHTS];

	out vec4 frag_colour;


	// array of offset direction for sampling
	vec3 gridSamplingDisk[20] = vec3[]
	(
		vec3(1, 1,  1), vec3( 1, -1,  1), vec3(-1, -1,  1), vec3(-1, 1,  1), 
		vec3(1, 1, -1), vec3( 1, -1, -1), vec3(-1, -1, -1), vec3(-1, 1, -1),
		vec3(1, 1,  0), vec3( 1, -1,  0), vec3(-1, -1,  0), vec3(-1, 1,  0),
		vec3(1, 0,  1), vec3(-1,  0,  1), vec3( 1,  0, -1), vec3(-1, 0, -1),
		vec3(0, 1,  1), vec3( 0, -1,  1), vec3( 0, -1, -1), vec3( 0, 1, -1)
	);

	vec3 CalcPointLight(PointLight light, vec3 normal, vec3 fragPos, vec3 viewDir, float shadow) 
	{
		vec3 lightDir = normalize(light.position - fragPos);
		// diffuse shading
		float diff = max(dot(lightDir, normal), 0.0);
		// specular shading
		vec3 reflectDir = reflect(lightDir, normal);
		float spec = pow(max(dot(viewDir, reflectDir), 0.0), nVal);
		// attenuation
		float distance    = length(light.position - fragPos);
		float attenuation = light.strength / (light.constant + light.linear * distance + 
					light.quadratic * (distance * distance));    
		// combine results
		//vec3 ambient  = light.color * ambientVal * diffuseVal;
		vec3 ambient = 0.3 * light.color * diffuseVal;
		vec3 diffuse  = light.color  * diff * diffuseVal;
		vec3 specular = vec3(0,0,0);

		if (diff < 0.0f) {
			specular = light.color * specularVal * spec;
			specular *= attenuation;
		}
		
		ambient  *= attenuation;
		diffuse  *= attenuation;
		
		return (ambient + diffuse + specular) * (1.0 - shadow);
	}

	float VectorToDepth (vec3 Vec)
	{
		vec3 AbsVec = abs(Vec);
		float LocalZcomp = max(AbsVec.x, max(AbsVec.y, AbsVec.z));

		// Replace f and n with the far and near plane values you used when
		//   you drew your cube map.
		const float f = 25.0;
		const float n = 1.0;

		float NormZComp = (f+n) / (f-n) - (2*f*n)/(f-n)/LocalZcomp;
		return (NormZComp + 1.0) * 0.5;
	}

	float ShadowCalculation(vec3 fragPos, PointLight light)
	{
		vec3 fragToLight = fragPos - light.position;
		float far_plane = 25.0;

		// float shadowVec = texture(depthMap, fragToLight).r;
		// // if (shadowVec + 0.0001 > VectorToDepth(fragToLight)) {
		// // 	return 1.0;
		// // }

		// // return 0.0;

		// shadowVec *= 25.0;
		float currentDepth = length(fragToLight);

		// float bias = 0.05;
		// float shadow = currentDepth - bias > shadowVec ? 1.0: 0.0;

		// //return shadowVec;
		// return shadow;

		float shadow = 0.0;
		float bias = 0.15;
		int samples = 20;
		float viewDistance = length(oCamPosition - fragPos);
		float diskRadius = (1.0 + (viewDistance / far_plane)) /25.0;

		for(int i = 0; i < samples; ++i)
		{
			float closestDepth = texture(depthMap, fragToLight + gridSamplingDisk[i] * diskRadius).r;
			closestDepth *= far_plane;   // undo mapping [0;1]
			if(currentDepth - bias > closestDepth)
				shadow += 1.0;
		}
		shadow /= float(samples);

		return shadow;
	}

	void main() {
		vec3 normal = normalize(normalInterp);
		vec3 result = vec3(0,0,0);
		vec3 viewDir = normalize(oCamPosition - oFragPosition);
		//float testShad = ShadowCalculation(oFragPosition, pointLights[0]);

		for (int i = 0; i < numLights; i++) {
			float shadow = ShadowCalculation(oFragPosition, pointLights[i]);
			result += CalcPointLight(pointLights[i], normal, oFragPosition, viewDir, shadow);
		}

		if (Alpha < 1.0) {
			frag_colour = vec4(result, Alpha);
		} else {
			frag_colour = vec4(result, Alpha);
		}

		//frag_colour = vec4(vec3(texture(depthMap, (oFragPosition - pointLights[0].position)).r), 1.0);
		//frag_colour = texture(depthMap, (oFragPosition - pointLights[0].position));
		//frag_colour = vec4(vec3(testShad), 1.0);
		//frag_colour = vec4(vec3(testShad), 1.0);

	}
	` + "\x00"
}