package shader

type BasicShader struct {
	fragShader string
	vertShader string
}

func (s BasicShader) GetFragShader() string {
	return s.fragShader
}

func (s BasicShader) GetVertShader() string {
	return s.vertShader
}

func (s *BasicShader) Setup() {
	s.vertShader = `
	#version 300 es
	//needed to add layout location for mac to work properly
	layout (location = 0) in vec3 aPosition;

	void main() {
		gl_Position = vec4(aPosition, 1.0);
	}
` + "\x00"

	s.fragShader = `
	#version 300 es
	precision highp float;

	out vec4 frag_colour;

	void main() {
		frag_colour = vec4(0.2, 0.3, 0.5, 1.0); 
	}
` + "\x00"
}
