package log

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/SyntSugar/ss-infra-go/consts"
	prome "github.com/SyntSugar/ss-infra-go/prometheus"
)

const (
	namespace                   = "infra"
	subsystem                   = "zap"
	ctxLogger consts.ContextKey = "ctx_logger"
)

var (
	globalLogger *Logger
	globalMetric metric
)

var ErrNil = errors.New("the value is null")

func NewLogger(loglevel, encoding string, samplingConfig *zap.SamplingConfig, fileLogPath string) (*Logger, error) {
	if fileLogPath == "" {
		fileLogPath = "stdout"
	}

	zapLogger, err := newZap(loglevel, encoding, samplingConfig, fileLogPath)
	if err != nil {
		return nil, err
	}

	debugger, err := newZap(zap.DebugLevel.String(), encoding, samplingConfig, fileLogPath)
	if err != nil {
		return nil, err
	}

	logger := &Logger{
		logger:   zapLogger,
		debugger: debugger,
	}
	logger = logger.WithOptions(zap.AddCallerSkip(1))

	return logger, nil
}

// DefaultSamplingConfig return the default sampling config for the zap logger.
// it would start sampling after exceeding initial entries(default was 64) per second
// and only sample the Thereafter(default was 10 which means sampling 10%) th entry.
func DefaultSamplingConfig() *zap.SamplingConfig {
	return &zap.SamplingConfig{
		Initial:    64, // log 64 entries per second before sampling
		Thereafter: 10, // log at 10th entry after exceeding the 64 entries
		Hook:       logSamplingMetrics,
	}
}

func GlobalLogger() *Logger {
	return globalLogger
}

type Logger struct {
	logger   *zap.Logger
	debugger *zap.Logger
}

// Named adds a new path segment to the logger's name. Segments are joined by
// periods. By default, Loggers are unnamed.
func (l *Logger) Named(s string) *Logger {
	if s == "" {
		return l
	}

	l.logger = l.logger.Named(s)
	l.debugger = l.debugger.Named(s)

	return l
}

// WithOptions clones the current Logger, applies the supplied Options, and
// returns the resulting Logger. It's safe to use concurrently.
func (l *Logger) WithOptions(opts ...zap.Option) *Logger {
	copy := l.clone()

	copy.logger = copy.logger.WithOptions(opts...)
	copy.debugger = copy.debugger.WithOptions(opts...)

	return copy
}

// With creates a child logger and adds structured context to it. Fields added
// to the child don't affect the parent, and vice versa.
func (l *Logger) With(fields ...zap.Field) *Logger {
	if len(fields) == 0 {
		return l
	}

	copy := l.clone()

	copy.logger = copy.logger.With(fields...)
	copy.debugger = copy.debugger.With(fields...)

	return copy
}

// Check returns a CheckedEntry if logging a message at the specified level
// is enabled. It's a completely optional optimization; in high-performance
// applications, Check can help avoid allocating a slice to hold fields.
func (l *Logger) Check(lvl zapcore.Level, msg string) *zapcore.CheckedEntry {
	return l.logger.Check(lvl, msg)
}

// GetZapLogger would return the zap logger
func (l *Logger) GetZapLogger() *zap.Logger {
	return l.logger
}

// Debug logs a message at DebugLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.logger.Debug(msg, fields...)
}

// DebugCtx logs a message at DebugLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (l *Logger) DebugCtx(ctx context.Context, msg string, fields ...zap.Field) {
	logger := l.logger
	if l.IsDynamicDebugEnabled(ctx) {
		logger = l.debugger
	}
	logger.Debug(msg, append(fields, GetContextFields(ctx)...)...)
}

// Info logs a message at InfoLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.logger.Info(msg, fields...)
}

// InfoCtx logs a message at InfoLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (l *Logger) InfoCtx(ctx context.Context, msg string, fields ...zap.Field) {
	fields = append(fields, GetContextFields(ctx)...)
	l.logger.Info(msg, fields...)
}

// Warn logs a message at WarnLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.logger.Warn(msg, fields...)
}

// WarnCtx logs a message at WarnLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (l *Logger) WarnCtx(ctx context.Context, msg string, fields ...zap.Field) {
	fields = append(fields, GetContextFields(ctx)...)
	l.logger.Warn(msg, fields...)
}

// Warn logs a message at WarnLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.logger.Error(msg, fields...)
}

// ErrorCtx logs a message at ErrorLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (l *Logger) ErrorCtx(ctx context.Context, msg string, fields ...zap.Field) {
	fields = append(fields, GetContextFields(ctx)...)
	l.logger.Error(msg, fields...)
}

// DPanic logs a message at DPanicLevel. The message includes any fields
// passed at the log site, as well as any fields accumulated on the logger.
//
// If the logger is in development mode, it then panics (DPanic means
// "development panic"). This is useful for catching errors that are
// recoverable, but shouldn't ever happen.
func (l *Logger) DPanic(msg string, fields ...zap.Field) {
	l.logger.DPanic(msg, fields...)
}

// DPanicCtx logs a message at DPanicLevel. The message includes any fields
// passed at the log site, as well as any fields accumulated on the logger.
//
// If the logger is in development mode, it then panics (DPanic means
// "development panic"). This is useful for catching errors that are
// recoverable, but shouldn't ever happen.
func (l *Logger) DPanicCtx(ctx context.Context, msg string, fields ...zap.Field) {
	fields = append(fields, GetContextFields(ctx)...)
	l.logger.DPanic(msg, fields...)
}

// Panic logs a message at PanicLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then panics, even if logging at PanicLevel is disabled.
func (l *Logger) Panic(msg string, fields ...zap.Field) {
	l.logger.DPanic(msg, fields...)
}

// PanicCtx logs a message at PanicLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then panics, even if logging at PanicLevel is disabled.
func (l *Logger) PanicCtx(ctx context.Context, msg string, fields ...zap.Field) {
	fields = append(fields, GetContextFields(ctx)...)
	l.logger.DPanic(msg, fields...)
}

// Fatal logs a message at FatalLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then calls os.Exit(1), even if logging at FatalLevel is
// disabled.
func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	l.logger.DPanic(msg, fields...)
}

// FatalCtx logs a message at FatalLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then calls os.Exit(1), even if logging at FatalLevel is
// disabled.
func (l *Logger) FatalCtx(ctx context.Context, msg string, fields ...zap.Field) {
	fields = append(fields, GetContextFields(ctx)...)
	l.logger.DPanic(msg, fields...)
}

// Sync calls the underlying Core's Sync method, flushing any buffered log
// entries. Applications should take care to call Sync before exiting.
func (l *Logger) Sync() error {
	return l.logger.Sync()
}

// Core returns the Logger's underlying zapcore.Core.
func (l *Logger) Core() zapcore.Core {
	return l.logger.Core()
}

func (l *Logger) clone() *Logger {
	copy := *l
	return &copy
}

// IsDynamicDebugEnabled check whether debug level logging enabled.
func (l *Logger) IsDynamicDebugEnabled(ctx context.Context) bool {
	if ctx == nil || l.Core().Enabled(zap.DebugLevel) {
		return false
	}
	enabled, ok := ctx.Value(consts.ContextKeyEnableDebugLogging).(bool)
	return ok && enabled
}

func init() {
	var err error
	globalLogger, err = NewLogger("info", "json", DefaultSamplingConfig(), "")
	if err != nil {
		panic("Failed to create global logger" + err.Error())
	}

	labels := []string{"level", "decision"}
	globalMetric.samplingCounter = prome.NewCounterHelper(namespace, subsystem, "sampling", labels...)
}

func newZap(loglevel, encoding string, samplingConfig *zap.SamplingConfig, logFilePath string) (*zap.Logger, error) {
	// Set default encoding if empty
	if encoding == "" {
		encoding = "json"
	}

	// Use default sampling config if not provided
	if samplingConfig == nil {
		samplingConfig = DefaultSamplingConfig()
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "event_time"
	encoderConfig.LevelKey = "severity"
	encoderConfig.MessageKey = "message"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{logFilePath}
	cfg.EncoderConfig = encoderConfig
	cfg.Encoding = encoding
	cfg.Sampling = samplingConfig
	var level zapcore.Level
	if err := level.Set(loglevel); err != nil {
		return nil, err
	}
	cfg.Level = zap.NewAtomicLevelAt(level)

	return cfg.Build(zap.AddStacktrace(zapcore.DPanicLevel), zap.AddCallerSkip(0))
}
