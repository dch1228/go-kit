package db

import (
	"time"

	glog "gorm.io/gorm/logger"
)

type Config struct {
	Host                  string        `mapstructure:"host"`
	Username              string        `mapstructure:"username"`
	Password              string        `mapstructure:"password"`
	Database              string        `mapstructure:"database"`
	MaxIdleConnections    int           `mapstructure:"max-idle-connections"`
	MaxOpenConnections    int           `mapstructure:"max-open-connections"`
	MaxConnectionLifeTime time.Duration `mapstructure:"max-connection-life-time"`

	Log LogConfig `mapstructure:"log"`
}

type LogConfig struct {
	Level                     glog.LogLevel `mapstructure:"level"`
	SlowThreshold             time.Duration `mapstructure:"slow-threshold"`
	IgnoreRecordNotFoundError bool          `mapstructure:"ignore-record-not-found-error"`
}
