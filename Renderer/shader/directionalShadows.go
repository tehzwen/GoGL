package shader

type DirectionalShadow struct {
	fragShader string
	vertShader string
	geoShader  string
}

func (s DirectionalShadow) GetFragShader() string {
	return s.fragShader
}

func (s DirectionalShadow) GetVertShader() string {
	return s.vertShader
}

func (s DirectionalShadow) GetGeometryShader() string {
	return s.geoShader
}

func (s *DirectionalShadow) Setup() {
	s.vertShader = `
	#version 410
	//needed to add layout location for mac to work properly
	layout (location = 0) in vec3 aPosition;

	uniform mat4 lightSpaceMatrix;
	uniform mat4 uModelMatrix;

	void main() {
		gl_Position = lightSpaceMatrix * uModelMatrix * vec4(aPosition, 1.0);
	}
` + "\x00"
	s.geoShader = ""
	s.fragShader = `
	#version 410
	precision highp float;

	void main() {
		gl_FragDepth = gl_FragCoord.z;
	}
` + "\x00"
}
