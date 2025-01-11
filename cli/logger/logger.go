package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Configure(verbose bool) {
	conf := zap.NewDevelopmentConfig()
	conf.DisableCaller = !verbose
	conf.DisableStacktrace = !verbose
	if !verbose {
		conf.EncoderConfig.TimeKey = ""
		conf.Level.SetLevel(zapcore.InfoLevel)
	}

	conf.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	zap.ReplaceGlobals(zap.Must(conf.Build()))
}
