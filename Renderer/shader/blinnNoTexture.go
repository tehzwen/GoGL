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

	struct DirectionalLight {
		float strength;
		vec3 color;
		vec3 direction;
		vec3 position;
		mat4 lightSpaceMatrix;
		sampler2D depthMap;
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
	uniform int numPointLights;
	uniform int numDirLights;
	uniform PointLight pointLights[MAX_LIGHTS];
	uniform DirectionalLight dirLight;

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

	float PointShadowCalculation(vec3 fragPos, PointLight light)
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

	float DirShadowCalculation(DirectionalLight light)
	{
		vec4 fragPosLightSpace = vec4(oFragPosition, 1.0) * light.lightSpaceMatrix;
		vec3 pos = fragPosLightSpace.xyz * 0.5 + 0.5;
		//vec3 projCoords = fragPosLightSpace.xyz / fragPosLightSpace.w;
		//projCoords = projCoords * 0.5 + 0.5; 
		//float currentDepth = projCoords.z;
		//float shadow = currentDepth > closestDepth  ? 1.0 : 0.0;
		//return shadow;
		float closestDepth = texture(light.depthMap, pos.xy).r;
		return closestDepth < pos.z ? 0.0 : 1.0;
		
	}

	vec3 CalcDirLight(DirectionalLight light, vec3 normal, vec3 viewDir)
	{
		float shadow = DirShadowCalculation(light);
		vec3 lightDir = normalize(light.direction);
		float diff = max(dot(normal, lightDir), 0.0);
		vec3 reflectDir = reflect(-lightDir, normal);
		float spec = pow(max(dot(viewDir, reflectDir), 0.0), nVal);
		vec3 ambient = light.color * ambientVal * diffuseVal;
		vec3 diffuse = light.color * diff * diffuseVal;
		vec3 specular = light.color * specularVal * spec;
		return (shadow * (diffuse + specular) + ambient);
	}

	vec3 CalcPointLight(PointLight light, vec3 normal, vec3 fragPos, vec3 viewDir) 
	{
		float shadow = PointShadowCalculation(oFragPosition, light);
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


		specular = light.color * specularVal * spec;
		//specular *= attenuation;

		
		ambient  *= attenuation;
		//diffuse  *= attenuation;
		
		return (ambient + (1.0 - shadow) * (diffuse + specular));
	}

	void main() {
		vec3 normal = normalize(normalInterp);
		vec3 result = vec3(0,0,0);
		vec3 viewDir = normalize(oCamPosition - oFragPosition);

		for (int i = 0; i < numPointLights; i++) {
			result += CalcPointLight(pointLights[i], normal, oFragPosition, viewDir);
		}
		//result += CalcDirLight(dirLight, normal, viewDir);

		if (Alpha < 1.0) {
			frag_colour = vec4(result, Alpha);
		} else {
			frag_colour = vec4(result, Alpha);
		}
	}
	` + "\x00"
}
