package log

import (
	"testing"
)

func TestLog(t *testing.T) {
	defer func() { _ = Sync() }()

	cfg := Config{
		Level:  "debug",
		Format: "json",
	}
	lg, err := New(cfg)
	if err != nil {
		panic(err)
	}
	SetLogger(lg)

	Debug("Testing")
	Debugf("Testing %s", "str")

	Info("Testing")
	Infof("Testing %s", "str")

	Warn("Testing")
	Warnf("Testing %s", "str")

	Error("Testing")
	Errorf("Testing %s", "str")
}
