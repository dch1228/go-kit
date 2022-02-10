package db

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DB = gorm.DB

func New(cfg Config) (*DB, error) {
	dsn := fmt.Sprintf(`%s:%s@tcp(%s)/%s?charset=utf8&parseTime=%t&loc=%s`,
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Database,
		true,
		"Local")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: NewLogger(cfg.Log),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConnections)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConnections)
	sqlDB.SetConnMaxLifetime(cfg.MaxConnectionLifeTime)

	return db, nil
}

func MustNew(cfg Config) *DB {
	db, err := New(cfg)
	if err != nil {
		panic(err)
	}
	return db
}
