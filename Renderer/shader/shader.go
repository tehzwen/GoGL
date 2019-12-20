package shader

type Shader interface {
	GetFragShader() string
	GetVertShader() string
	Setup()
}
