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
	SetLogger(MustNew(cfg))

	Debug("Testing")
	Debugf("Testing %s", "str")

	Info("Testing")
	Infof("Testing %s", "str")

	Warn("Testing")
	Warnf("Testing %s", "str")

	Error("Testing")
	Errorf("Testing %s", "str")

	lg := With(String("str", "str"))
	lg.Info("info")

	lg = lg.With(Int("int", 1))
	lg.Info("info")
}
