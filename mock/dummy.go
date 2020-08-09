package mock

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger, _ = zap.Config{
	Encoding:    "json",
	Level:       zap.NewAtomicLevelAt(zapcore.FatalLevel),
	OutputPaths: []string{"stdout"},
}.Build()
