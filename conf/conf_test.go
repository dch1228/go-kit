package conf

import (
	"reflect"
	"testing"

	"gopkg.in/mcuadros/go-defaults.v1"
)

func TestLoad(t *testing.T) {
	v := struct {
	}{}

	defaults.SetDefaults(&v)

	objVal := reflect.ValueOf(v)
	srvCfgVal := objVal.FieldByName("ServerConfig")
	name = srvCfgVal.FieldByName("Env").String()

}
