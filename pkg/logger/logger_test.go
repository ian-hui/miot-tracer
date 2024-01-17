package logger

import (
	"testing"
)

func TestLogger(t *testing.T) {
	// 1. init logger
	Miotlogger.Info("init logger success")
	Miotlogger.Sync()
}
