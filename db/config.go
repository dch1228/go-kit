package db

import (
	"time"

	glog "gorm.io/gorm/logger"
)

type Config struct {
	Host                  string
	Username              string
	Password              string
	Database              string
	MaxIdleConnections    int           `default:"10"`
	MaxOpenConnections    int           `default:"10"`
	MaxConnectionLifeTime time.Duration `default:"3m"`

	Log LogConfig `mapstructure:"log"`
}

type LogConfig struct {
	Level                     glog.LogLevel `default:"4"`
	SlowThreshold             time.Duration `default:"1s"`
	IgnoreRecordNotFoundError bool          `default:"true"`
}
