package logger

import (
	"fmt"
	"io"
	"log"

	"github.com/GoTaskFlow/pkg/logger/model"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLogger struct {
	entry     *zap.Logger
	zapConfig zap.Config
	cfg       *model.Config
}

func NewZapLogger(config *model.Config) (*ZapLogger, error) {
	return newZapLogger(config)

}

func newZapLogger(config *model.Config) (*ZapLogger, error) {
	zapConfig := zap.NewProductionConfig()

	switch config.Level {
	case model.DebugLevel, model.TraceLevel:
		zapConfig.Development = true
	default:
		zapConfig.Development = false
	}

	zapConfig.Level = zap.NewAtomicLevelAt(config.Level.ZapLevel())
	zapConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	zapConfig.EncoderConfig.CallerKey = ""
	log := &ZapLogger{zapConfig: zapConfig, cfg: config}

	if err := log.build(); err != nil {
		return nil, err
	}

	return log, nil

}
func (log *ZapLogger) build() error {
	l, err := log.zapConfig.Build()
	if err != nil {
		return fmt.Errorf("failed to build logger: %w", err)
	}
	// if we want logs in a file
	// log.cfg.Output, err = os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }
	if log.cfg.Output != nil {
		l = log.writerLog(log.cfg.Output)
	}

	log.entry = l
	return nil

}
func (log *ZapLogger) writerLog(output io.Writer) *zap.Logger {
	encoder := zapcore.NewJSONEncoder(log.zapConfig.EncoderConfig)

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(output), log.zapConfig.Level),
	)

	return zap.New(core)
}

func (log *ZapLogger) Trace(msg string) {
	log.entry.Debug(msg)
}

func (log *ZapLogger) Tracef(format string, args ...interface{}) {
	log.entry.Debug(fmt.Sprintf(format, args...))
}

func (log *ZapLogger) Debug(msg string) {
	log.entry.Debug(msg)
}

func (log *ZapLogger) Debugf(format string, args ...interface{}) {
	log.entry.Debug(fmt.Sprintf(format, args...))
}

func (log *ZapLogger) Info(msg string) {
	log.entry.Info(msg)
}

func (log *ZapLogger) Infof(format string, args ...interface{}) {
	log.entry.Info(fmt.Sprintf(format, args...))
}

func (log *ZapLogger) Warn(msg string) {
	log.entry.Warn(msg)
}

func (log *ZapLogger) Warnf(format string, args ...interface{}) {
	log.entry.Warn(fmt.Sprintf(format, args...))
}

func (log *ZapLogger) Error(msg string) {
	log.entry.Error(msg)
}

func (log *ZapLogger) Errorf(format string, args ...interface{}) {
	log.entry.Error(fmt.Sprintf(format, args...))
}

func (log *ZapLogger) Fatal(msg string) {
	log.entry.Fatal(msg)
}

func (log *ZapLogger) Fatalf(format string, args ...interface{}) {
	log.entry.Fatal(fmt.Sprintf(format, args...))
}

//nolint:ireturn,gofmt // implements model.Logger interface
func (log *ZapLogger) WithField(key string, value interface{}) model.Logger {
	return &ZapLogger{entry: log.entry.With(zap.Any(key, value))}
}

//nolint:ireturn,gofmt // implements model.Logger interface
func (log *ZapLogger) WithFields(fields model.Fields) model.Logger {
	zFields := make([]zapcore.Field, 0)
	for key, value := range fields {
		zFields = append(zFields, zap.Any(key, value))
	}
	return &ZapLogger{entry: log.entry.With(zFields...)}
}

//nolint:ireturn // implements model.Logger interface
func (log *ZapLogger) WithError(err error) model.Logger {
	return &ZapLogger{entry: log.entry.With(zap.Error(err))}
}

// ToStdLogger creates a logger that matches std library.
func (log *ZapLogger) ToStdLogger() *log.Logger {
	return zap.NewStdLog(log.entry)
}

func (log *ZapLogger) SetLevel(lvl model.Level) error {
	switch lvl {
	case model.DebugLevel, model.TraceLevel:
		log.zapConfig.Development = true
	default:
		log.zapConfig.Development = false
	}

	log.zapConfig.Level = zap.NewAtomicLevelAt(lvl.ZapLevel())
	return log.build()
}
