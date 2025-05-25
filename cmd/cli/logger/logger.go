package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Configure() {
	conf := zap.NewDevelopmentConfig()
	level := zap.NewAtomicLevel()
	logLevel := os.Getenv("LOG_LEVEL")
	if err := level.UnmarshalText([]byte(logLevel)); err != nil {
		level.SetLevel(zapcore.WarnLevel) // default
	}
	isDebug := level.Level() == zap.DebugLevel
	if !isDebug {
		conf.DisableCaller = true
		conf.DisableStacktrace = true
		conf.EncoderConfig.TimeKey = ""
	}
	conf.Level.SetLevel(level.Level())
	// conf.EncoderConfig.EncodeLevel = func(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	// 	if l != zapcore.InfoLevel || verbose {
	// 		zapcore.CapitalColorLevelEncoder(l, enc)
	// 	}
	// }
	conf.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	zap.ReplaceGlobals(zap.Must(conf.Build()))
}
