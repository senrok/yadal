package layers

import (
	"go.uber.org/zap"
)

func ExampleNewLoggerAdapter() {
	logger, _ := zap.NewProduction()
	s := logger.Sugar()
	_ = NewLoggerAdapter(s.Info, s.Infof)
}

func ExampleNewLoggingLayer() {
	// new logger
	logger, _ := zap.NewProduction()
	s := logger.Sugar()

	_ = NewLoggingLayer(SetLogger(NewLoggerAdapter(s.Info, s.Infof)))
}
