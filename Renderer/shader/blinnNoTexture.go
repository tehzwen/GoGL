package shader

type BlinnNoTexture struct {
	fragShader string
	vertShader string
}

func (s BlinnNoTexture) GetFragShader() string {
	return s.fragShader
}

func (s BlinnNoTexture) GetVertShader() string {
	return s.vertShader
}

func (s *BlinnNoTexture) Setup() {
	s.vertShader = `
	#version 300 es
	//needed to add layout location for mac to work properly
	layout (location = 0) in vec3 aPosition;
	layout (location = 1) in vec3 aNormal;

	out vec3 oNormal;
	out vec3 normalInterp;
	out vec3 oFragPosition;
	
	uniform mat4 uProjectionMatrix;
	uniform mat4 uViewMatrix;
	uniform mat4 uModelMatrix;

	void main() {
		mat4 normalMatrix = transpose(inverse(uModelMatrix));
		oNormal = normalize((uModelMatrix * vec4(aNormal, 1.0)).xyz);
		normalInterp = vec3(normalMatrix * vec4(aNormal, 0.0));
		oFragPosition = (uModelMatrix * vec4(aPosition, 1.0)).xyz;
		gl_Position = uProjectionMatrix * uViewMatrix * uModelMatrix * vec4(aPosition, 1.0); 
	}
` + "\x00"

	s.fragShader = `
	#version 300 es
	precision highp float;
	#define MAX_LIGHTS 128

	in vec3 oFragPosition;
	in vec3 normalInterp;
	in vec3 oNormal;

	uniform vec3 cameraPosition;
	uniform vec3 diffuseVal;
	uniform vec3 ambientVal;
	uniform vec3 specularVal;
	uniform float nVal;
	uniform int numLights;
	uniform vec3 lightPositions[MAX_LIGHTS];
	uniform vec3 lightColours[MAX_LIGHTS];
	uniform float lightStrengths[MAX_LIGHTS];

	out vec4 frag_colour;

	void main() {
		
		vec3 diffuse = vec3(0, 0, 0);
		vec3 ambient = vec3(0, 0, 0);
		vec3 specular = vec3(0, 0, 0);
		vec3 normal = normalize(normalInterp);

		for (int i = 0; i < numLights; i++) {
			vec3 nCameraPosition = normalize(cameraPosition); // Normalize the camera Position
			vec3 V = normalize(nCameraPosition - oFragPosition);

			vec3 lightDirection = normalize(lightPositions[i] - oFragPosition);
			float diff = max(dot(normal, lightDirection), 0.0);
			vec3 reflectDir = reflect(-lightDirection, normal);
			float spec = pow(max(dot(V, reflectDir), 0.0), nVal);
			float distance = length(lightPositions[i] - oFragPosition);
			float attenuation = 1.0 / (distance * distance);
			attenuation *= lightStrengths[i];

			ambient += ambientVal * lightColours[i] * diffuseVal;
			diffuse += diffuseVal * lightColours[i] * diff;

			if (diff > 0.0f) {
				specular += specularVal * lightColours[i] * spec;
				specular *= attenuation;
			}
			ambient *= attenuation; //causes much darker scene
			diffuse *= attenuation;
		}
		frag_colour = vec4(diffuse + ambient + specular, 1.0); 
	}
` + "\x00"
}
