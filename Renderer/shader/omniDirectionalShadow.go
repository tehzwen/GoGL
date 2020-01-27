package shader

type OmniDirectionalShadow struct {
	fragShader string
	vertShader string
	geoShader  string
}

func (s OmniDirectionalShadow) GetFragShader() string {
	return s.fragShader
}

func (s OmniDirectionalShadow) GetVertShader() string {
	return s.vertShader
}

func (s OmniDirectionalShadow) GetGeometryShader() string {
	return s.geoShader
}

func (s *OmniDirectionalShadow) Setup() {
	s.vertShader = `
	#version 410
	//needed to add layout location for mac to work properly
	layout (location = 0) in vec3 aPosition;
	

	uniform mat4 uModelMatrix;

	void main() {
		gl_Position = uModelMatrix * vec4(aPosition, 1.0);
	}
` + "\x00"
	s.geoShader = `
	#version 410
	layout (triangles) in;
	layout (triangle_strip, max_vertices=18) out;

	uniform mat4 shadowMatrices[6];

	out vec4 FragPos; // FragPos from GS (output per emitvertex)

	void main()
	{
		for(int face = 0; face < 6; ++face)
		{
			gl_Layer = face; // built-in variable that specifies to which face we render.
			for(int i = 0; i < 3; ++i) // for each triangle's vertices
			{
				FragPos = gl_in[i].gl_Position;
				gl_Position = shadowMatrices[face] * FragPos;
				EmitVertex();
			}    
			EndPrimitive();
		}
		
	}  
	`
	s.fragShader = `
	#version 410
	precision highp float;

	in vec4 FragPos;

	uniform vec3 lightPos;

	void main() {
		float far_plane = 25.0f;

		//try with the first light in the array
		float lightDistance = length(lightPos - FragPos.xyz);
		lightDistance = lightDistance / far_plane;
		gl_FragDepth = lightDistance;
	}
` + "\x00"
}
