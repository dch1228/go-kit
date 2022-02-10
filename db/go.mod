module github.com/dch1228/go-kit/db

go 1.17

require (
	github.com/dch1228/go-kit/log v0.0.0-00010101000000-000000000000
	gorm.io/driver/mysql v1.2.3
	gorm.io/gorm v1.22.5
)

require (
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.4 // indirect
	go.opentelemetry.io/otel v1.3.0 // indirect
	go.opentelemetry.io/otel/trace v1.3.0 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.21.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
)

replace github.com/dch1228/go-kit/log => ../log
