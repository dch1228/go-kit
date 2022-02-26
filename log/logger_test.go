package log

import (
	"errors"
	"testing"
)

func TestLog(t *testing.T) {
	lg := New(Config{
		Level: "info",
	})
	defer func() { _ = lg.Sync() }()

	SetLogger(lg)

	Info("Info")

	Warn("Warn")

	Error("Error", errors.New("test error"))

	lg = With(String("str", "str"))
	lg.Info("info")

	lg = lg.With(Int("int", 1))
	lg.Info("info")
}
