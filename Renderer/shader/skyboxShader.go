package shader

type SkyboxShader struct {
	fragShader string
	vertShader string
	geoShader  string
}

func (s SkyboxShader) GetFragShader() string {
	return s.fragShader
}

func (s SkyboxShader) GetVertShader() string {
	return s.vertShader
}

func (s SkyboxShader) GetGeometryShader() string {
	return s.geoShader
}

func (s *SkyboxShader) Setup() {
	s.vertShader = `
	#version 410
	//needed to add layout location for mac to work properly
	layout (location = 0) in vec3 aPosition;

	uniform mat4 uProjectionMatrix;
	uniform mat4 uViewMatrix;

	out vec3 TexCoords;

	void main() {
		TexCoords = aPosition;
		vec4 pos = uProjectionMatrix * uViewMatrix * vec4(aPosition, 1.0);
		gl_Position = pos.xyww;
	}
` + "\x00"
	s.geoShader = ""
	s.fragShader = `
	#version 410
	precision highp float;

	in vec3 TexCoords;

	uniform samplerCube skybox;

	out vec4 frag_colour;

	void main() {
		frag_colour = texture(skybox, TexCoords);
	}
` + "\x00"
}
