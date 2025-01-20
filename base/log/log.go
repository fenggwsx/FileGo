package log

import (
	"github.com/fenggwsx/filego/base/conf"
	"go.uber.org/zap"
)

func newDevelopmentLogger() (*zap.Logger, error) {
	return zap.Config{
		Level:            zap.NewAtomicLevelAt(zap.DebugLevel),
		Development:      true,
		Encoding:         "console",
		EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}.Build(zap.AddStacktrace(zap.WarnLevel))
}

func newProductionLogger() (*zap.Logger, error) {
	return zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}.Build(zap.AddStacktrace(zap.ErrorLevel))
}

func NewLogger(config *conf.Config) (*zap.Logger, error) {
	if config.Development {
		return newDevelopmentLogger()
	} else {
		return newProductionLogger()
	}
}
