package db

import (
	"testing"
	"time"

	glog "gorm.io/gorm/logger"
)

func TestDB(t *testing.T) {
	db := MustNew(Config{
		Host:                  "127.0.0.1:3306",
		Username:              "root",
		Password:              "root",
		Database:              "test",
		MaxIdleConnections:    100,
		MaxOpenConnections:    100,
		MaxConnectionLifeTime: 10 * time.Minute,
		Log: LogConfig{
			Level:                     glog.Info,
			SlowThreshold:             1 * time.Second,
			IgnoreRecordNotFoundError: true,
		},
	})

	out := make(map[string]interface{})
	err := db.Raw("select sleep(1)").Scan(&out).Error
	if err != nil {
		t.Fatal(err)
	}

	db.Clauses()
}
