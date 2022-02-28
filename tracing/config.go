package tracing

type Config struct {
	Name         string  `mapstructure:"name"`
	Endpoint     string  `mapstructure:"endpoint"`
	SamplerRatio float64 `mapstructure:"sampler-ratio"`
}
