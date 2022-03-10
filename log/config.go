package log

type Config struct {
	Level  string `validate:"oneof=info warn error" default:"info"`
	Format string `validate:"oneof=text json" default:"text"`

	EnableTrace bool `default:"true"`

	File FileConfig
}

type FileConfig struct {
	Filename   string
	MaxSize    int
	MaxAge     int
	MaxBackups int
}
