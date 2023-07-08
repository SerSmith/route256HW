package logger

import (
	"context"
	"io"
	"os"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	global	   *zap.SugaredLogger
	defaultLevel = zap.NewAtomicLevelAt(zap.InfoLevel)
)

func init() {
	global = New(defaultLevel, os.Stdout)
}

func New(level zap.AtomicLevel, sink io.Writer, opts ...zap.Option) *zap.SugaredLogger {
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zapcore.EncoderConfig{
			TimeKey:		"ts",
			LevelKey:	   "level",
			NameKey:		"logger",
			CallerKey:	  "caller",
			MessageKey:	 "message",
			StacktraceKey:  "stacktrace",
			LineEnding:	 zapcore.DefaultLineEnding,
			EncodeLevel:	zapcore.LowercaseLevelEncoder,
			EncodeTime:	 zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}),
		zapcore.AddSync(sink),
		level,
	)

	return zap.New(core, opts...).Sugar()
}

// SetLoggerByEnvironment - not thread safe
func SetLoggerByEnvironment(environment string) {
	if environment == "PRODUCTION" {
		global = New(zap.NewAtomicLevelAt(zap.ErrorLevel), os.Stdout,
			zap.WithCaller(false),
			zap.AddStacktrace(zap.NewAtomicLevelAt(zap.PanicLevel)),
		)
	}
}

func Info(args ...interface{}) {
	global.Info(args...)
}

func Error(args ...interface{}) {
	global.Error(args...)
}

func Errorf(ctx context.Context, method, template string, args ...interface{}) {
	withTraceID(ctx).Desugar().
		With(zap.String("method", method)).Sugar().Errorf(template, args...)
}

func Fatal(args ...interface{}) {
	global.Fatal(args...)
}

func Fatalf(template string, args ...interface{}) {
	global.Fatalf(template, args...)
}

func withTraceID(ctx context.Context) *zap.SugaredLogger {
	span := opentracing.SpanFromContext(ctx)
	if span == nil {
		return global
	}

	if sc, ok := span.Context().(jaeger.SpanContext); ok {
		return global.Desugar().With(
			zap.Stringer("trace_id", sc.TraceID()),
			zap.Stringer("span_id", sc.SpanID()),
		).Sugar()
	}

	return global
}

