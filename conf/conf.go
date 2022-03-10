package conf

import (
	"context"
	"fmt"
	"io/ioutil"
	"sync"

	"github.com/go-playground/validator/v10"
	"gopkg.in/mcuadros/go-defaults.v1"
	"gopkg.in/yaml.v2"

	"github.com/dch1228/go-kit/log"
	"github.com/dch1228/go-kit/tracing"
)

var (
	once sync.Once
)

type ServerConfig struct {
	Name    string
	Env     string `validate:"oneof=local dev prod" default:"local"`
	Log     log.Config
	Tracing tracing.Config
}

func Load(path string, v interface{}) error {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(content, &v); err != nil {
		return err
	}

	defaults.SetDefaults(&v)

	if err := validator.New().Struct(v); err != nil {
		return err
	}

	return nil
}

func MustLoad(path string, v interface{}) {
	if err := Load(path, v); err != nil {
		panic(fmt.Errorf("error: config file %s, %s", path, err))
	}
}

func (cfg *ServerConfig) Setup() (err error) {
	once.Do(func() {
		if err = log.Setup(cfg.Log); err != nil {
			return
		}
		if err = tracing.Setup(cfg.Tracing); err != nil {
			return
		}
	})
	return
}

func (cfg *ServerConfig) MustSetup() {
	if err := cfg.Setup(); err != nil {
		panic(fmt.Errorf("error: %s", err))
	}
}

func (cfg *ServerConfig) Cleanup() {
	if err := tracing.Shutdown(context.Background()); err != nil {
		log.Error("tracing.Shutdown", err)
	}

	_ = log.Sync()
}
