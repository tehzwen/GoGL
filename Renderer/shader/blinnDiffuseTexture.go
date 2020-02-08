package shader

type BlinnDiffuseTexture struct {
	fragShader string
	vertShader string
	geoShader  string
}

func (s BlinnDiffuseTexture) GetFragShader() string {
	return s.fragShader
}

func (s BlinnDiffuseTexture) GetVertShader() string {
	return s.vertShader
}

func (s BlinnDiffuseTexture) GetGeometryShader() string {

	return s.geoShader
}

func (s *BlinnDiffuseTexture) Setup() {
	s.vertShader = `
	#version 410
	//needed to add layout location for mac to work properly
	layout (location = 0) in vec3 aPosition;
	layout (location = 1) in vec3 aNormal;
	layout (location = 2) in vec2 aUV;

	out vec3 oNormal;
	out vec3 normalInterp;
	out vec3 oFragPosition;
	out vec2 oUV;
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
		oUV = aUV;
		oCamPosition =  (uViewMatrix * vec4(cameraPosition, 1.0)).xyz;
		gl_Position = uProjectionMatrix * uViewMatrix * uModelMatrix * vec4(aPosition, 1.0); 
	}
` + "\x00"

	s.geoShader = ""

	s.fragShader = `
	#version 410
	precision highp float;
	#define MAX_LIGHTS 20

	struct PointLight {
		vec3 position;
		float strength;
		float constant;
		float linear;
		float quadratic; 
		vec3 color;
		samplerCube depthMap;
	};

	in vec3 oFragPosition;
	in vec3 normalInterp;
	in vec3 oNormal;
	in vec3 oCamPosition;
	in vec2 oUV;

	uniform vec3 diffuseVal;
	uniform vec3 ambientVal;
	uniform vec3 specularVal;
	uniform float nVal;
	uniform float Alpha;
	uniform int numLights;
	uniform sampler2D uDiffuseTexture;
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

	float ShadowCalculation(vec3 fragPos, PointLight light)
	{
		vec3 fragToLight = fragPos - light.position;
		float far_plane = 25.0;
		float currentDepth = length(fragToLight);

		float shadow = 0.0;
		float bias = 0.15;
		int samples = 20;
		float viewDistance = length(oCamPosition - fragPos);
		float diskRadius = (1.0 + (viewDistance / far_plane)) /25.0;

		for(int i = 0; i < samples; ++i)
		{
			float closestDepth = texture(light.depthMap, fragToLight + gridSamplingDisk[i] * diskRadius).r;
			closestDepth *= far_plane;   // undo mapping [0;1]
			if(currentDepth - bias > closestDepth)
				shadow += 1.0;
		}
		shadow /= float(samples);

		return shadow;
	}

	vec3 CalcPointLight(PointLight light, vec3 normal, vec3 fragPos, vec3 viewDir, vec3 textureVal) 
	{
		float shadow = ShadowCalculation(oFragPosition, light);
		vec3 lightDir = normalize(light.position - fragPos);
		// diffuse shading
		float diff = max(dot(normal, lightDir), 0.0);
		// specular shading
		vec3 reflectDir = reflect(lightDir, normal);
		float spec = pow(max(dot(viewDir, reflectDir), 0.0), nVal);
		// attenuation
		float distance    = length(light.position - fragPos);
		float attenuation = light.strength / (light.constant + light.linear * distance + 
					   light.quadratic * (distance * distance));    
		// combine results
		vec3 ambient  = light.color * ambientVal * diffuseVal * textureVal;
		vec3 diffuse  = light.color  * diff * diffuseVal * textureVal;
		vec3 specular = vec3(0,0,0);

		if (diff < 0.0f) {
			specular = light.color * specularVal * spec * textureVal;
			specular *= attenuation;
		}
		
		ambient  *= attenuation;
		diffuse  *= attenuation;
		
		return (ambient + (1.0 - shadow) * (diffuse + specular));
	}

	void main() {
		vec3 normal = normalize(normalInterp);
		vec3 result = vec3(0,0,0);
		vec3 viewDir = normalize(oCamPosition - oFragPosition);

		vec4 texColor = texture(uDiffuseTexture, oUV);

		for (int i = 0; i < numLights; i++) {
			result += CalcPointLight(pointLights[i], normal, oFragPosition, viewDir, texColor.xyz);
		}

		if (texColor.w < 0.1) {
			discard;
		}

		frag_colour = vec4(result, Alpha);
		//frag_colour = vec4(0.5, 0.0, 0.0, 1.0);
	}
` + "\x00"
}
