package log

import (
	"errors"
	"testing"
)

func TestLog(t *testing.T) {
	MustSetup(Config{
		Level:       "info",
		EnableTrace: true,
	})

	Info("Info")

	Warn("Warn")

	Error("Error", errors.New("test error"))

	lg := With(String("str", "str"))
	lg.Info("info")

	lg = lg.With(Int("int", 1))
	lg.Info("info")
}
