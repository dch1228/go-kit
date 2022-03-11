package tracing

type Config struct {
	Name         string
	Endpoint     string  `default:"http://127.0.0.1:14268/api/traces"`
	SamplerRatio float64 `default:"1"`
}
