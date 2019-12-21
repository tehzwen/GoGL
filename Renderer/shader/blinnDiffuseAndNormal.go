package shader

type BlinnDiffuseAndNormal struct {
	fragShader string
	vertShader string
}

func (s BlinnDiffuseAndNormal) GetFragShader() string {
	return s.fragShader
}

func (s BlinnDiffuseAndNormal) GetVertShader() string {
	return s.vertShader
}

func (s *BlinnDiffuseAndNormal) Setup() {
	s.vertShader = `
	#version 330 core
	//needed to add layout location for mac to work properly
	layout (location = 0) in vec3 aPosition;
	layout (location = 1) in vec3 aNormal;
	layout (location = 2) in vec2 aUV;

	out vec3 oNormal;
	out vec3 normalInterp;
	out vec3 oFragPosition;
	out vec2 oUV;
	
	uniform mat4 uProjectionMatrix;
	uniform mat4 uViewMatrix;
	uniform mat4 uModelMatrix;

	void main() {
		mat4 normalMatrix = transpose(inverse(uModelMatrix));
		oNormal = normalize((uModelMatrix * vec4(aNormal, 1.0)).xyz);
		normalInterp = vec3(normalMatrix * vec4(aNormal, 0.0));
		oFragPosition = (uModelMatrix * vec4(aPosition, 1.0)).xyz;
		oUV = aUV;
		gl_Position = uProjectionMatrix * uViewMatrix * uModelMatrix * vec4(aPosition, 1.0); 
	}
` + "\x00"

	s.fragShader = `
	#version 410
	precision highp float;
	#define MAX_LIGHTS 128

	struct PointLight {
		vec3 position;
		float strength;
		vec3 color;
	};

	in vec3 oFragPosition;
	in vec3 normalInterp;
	in vec3 oNormal;
	in vec2 oUV;

	uniform vec3 cameraPosition;
	uniform vec3 diffuseVal;
	uniform vec3 ambientVal;
	uniform vec3 specularVal;
	uniform float nVal;
	uniform int numLights;
	uniform sampler2D uDiffuseTexture;
	uniform PointLight pointLights[MAX_LIGHTS];

	out vec4 frag_colour;

	void main() {
		
		vec3 diffuse = vec3(0, 0, 0);
		vec3 ambient = vec3(0, 0, 0);
		vec3 specular = vec3(0, 0, 0);
		vec3 normal = normalize(normalInterp);

		for (int i = 0; i < numLights; i++) {
			vec3 nCameraPosition = normalize(cameraPosition); // Normalize the camera Position
			vec3 V = normalize(nCameraPosition - oFragPosition);

			vec3 lightDirection = normalize(pointLights[i].position - oFragPosition);
			float diff = max(dot(normal, lightDirection), 0.0);
			vec3 reflectDir = reflect(-lightDirection, normal);
			float spec = pow(max(dot(V, reflectDir), 0.0), nVal);
			float distance = length(pointLights[i].position - oFragPosition);
			float attenuation = 1.0 / (distance * distance);
			//attenuation *= lightStrengths[i];

			ambient += ambientVal * pointLights[i].color * diffuseVal;
			diffuse += (diffuseVal * pointLights[i].color * diff) + pointLights[i].strength;

			if (diff > 0.0f) {
				specular += specularVal * pointLights[i].color * spec;
				specular *= attenuation;
			}
			ambient *= attenuation; //causes much darker scene
			diffuse *= attenuation;
		}

		vec4 textureColor = texture(uDiffuseTexture, oUV);
		
		frag_colour = vec4((diffuse + ambient + specular) * textureColor.rgb, 1.0); 
	}
	` + "\x00"
}
