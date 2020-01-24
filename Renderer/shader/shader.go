package shader

type Shader interface {
	GetFragShader() string
	GetVertShader() string
	GetGeometryShader() (string, error)
	Setup()
}
