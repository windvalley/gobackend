package log

import (
	"context"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ctxKey int

const logContextKey ctxKey = iota

// Logger ...
type Logger struct {
	zapLogger *zap.Logger
}

// Create a new logger instance.
func newLogger(l *zap.Logger) *Logger {
	return &Logger{
		zapLogger: l,
	}
}

// WithName custom logger name.
func (l *Logger) WithName(name string) *Logger {
	logger := l.zapLogger.Named(name)
	return newLogger(logger)
}

// WithFields custom other log entry fileds.
func (l *Logger) WithFields(fields ...Field) *Logger {
	logger := l.zapLogger.With(fields...)
	return newLogger(logger)
}

// WithHooks is different from SetHooks,
// SetHooks is for global logger,
// WithHooks is for the new logger.
func (l *Logger) WithHooks(hooks ...Hook) *Logger {
	logger := l.zapLogger.WithOptions(zap.Hooks(hooks...))
	return newLogger(logger)
}

// ToContext put logger to context.
func (l *Logger) ToContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, logContextKey, l)
}

// FromContext return logger from context.
func (l *Logger) FromContext(ctx context.Context) *Logger {
	if ctx != nil {
		logger := ctx.Value(logContextKey)
		if logger != nil {
			return logger.(*Logger)
		}
	}

	return l.WithName("UnknownContext")
}

// C get logger from gin.Context.
//
// Usage example:
//
// This is a middleware that put logger into gin.Context:
//
//func Context() gin.HandlerFunc {
//return func(c *gin.Context) {
//l := log.WithFields(
//log.String("x-request-id", c.GetString(XRequestIDKey)),
//log.String("username", c.GetString(UsernameKey)),
//)
//c.Set(log.ContextLoggerName, l)
//c.Next()
//}
//}
//
// Get logger that with fileds from gin.Context:
//
//func (u *UserController) Get(c *gin.Context) {
//log.C(c).Debug("user get called")
//}
//
func (l *Logger) C(ctx context.Context) *Logger {
	ctxLogger, ok := ctx.(*gin.Context).Get(ContextLoggerName)
	if !ok {
		return l
	}

	cl, ok := ctxLogger.(*Logger)
	if !ok {
		return l
	}

	return cl
}

// Debug ...
func (l *Logger) Debug(msg string, fields ...Field) {
	l.zapLogger.Debug(msg, fields...)
}

// Debugf ...
func (l *Logger) Debugf(format string, v ...interface{}) {
	l.zapLogger.Sugar().Debugf(format, v...)
}

// Debugw ...
func (l *Logger) Debugw(msg string, keysAndValues ...interface{}) {
	l.zapLogger.Sugar().Debugw(msg, keysAndValues...)
}

// Info ...
func (l *Logger) Info(msg string, fields ...Field) {
	l.zapLogger.Info(msg, fields...)
}

// Infof ...
func (l *Logger) Infof(format string, v ...interface{}) {
	l.zapLogger.Sugar().Infof(format, v...)
}

// Infow ...
func (l *Logger) Infow(msg string, keysAndValues ...interface{}) {
	l.zapLogger.Sugar().Infow(msg, keysAndValues...)
}

// Warn ...
func (l *Logger) Warn(msg string, fields ...Field) {
	l.zapLogger.Warn(msg, fields...)
}

// Warnf ...
func (l *Logger) Warnf(format string, v ...interface{}) {
	l.zapLogger.Sugar().Warnf(format, v...)
}

// Warnw ...
func (l *Logger) Warnw(msg string, keysAndValues ...interface{}) {
	l.zapLogger.Sugar().Warnw(msg, keysAndValues...)
}

// Error ...
func (l *Logger) Error(msg string, fields ...Field) {
	l.zapLogger.Error(msg, fields...)
}

// Errorf ...
func (l *Logger) Errorf(format string, v ...interface{}) {
	l.zapLogger.Sugar().Errorf(format, v...)
}

// Errorw ...
func (l *Logger) Errorw(msg string, keysAndValues ...interface{}) {
	l.zapLogger.Sugar().Errorw(msg, keysAndValues...)
}

// DPanic ...
func (l *Logger) DPanic(msg string, fields ...Field) {
	l.zapLogger.DPanic(msg, fields...)
}

// DPanicf ...
func (l *Logger) DPanicf(format string, v ...interface{}) {
	l.zapLogger.Sugar().DPanicf(format, v...)
}

// DPanicw ...
func (l *Logger) DPanicw(msg string, keysAndValues ...interface{}) {
	l.zapLogger.Sugar().DPanicw(msg, keysAndValues...)
}

// Panic ...
func (l *Logger) Panic(msg string, fields ...Field) {
	l.zapLogger.Panic(msg, fields...)
}

// Panicf ...
func (l *Logger) Panicf(format string, v ...interface{}) {
	l.zapLogger.Sugar().Panicf(format, v...)
}

// Panicw ...
func (l *Logger) Panicw(msg string, keysAndValues ...interface{}) {
	l.zapLogger.Sugar().Panicw(msg, keysAndValues...)
}

// Fatal ...
func (l *Logger) Fatal(msg string, fields ...Field) {
	l.zapLogger.Fatal(msg, fields...)
}

// Fatalf ...
func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.zapLogger.Sugar().Fatalf(format, v...)
}

// Fatalw ...
func (l *Logger) Fatalw(msg string, keysAndValues ...interface{}) {
	l.zapLogger.Sugar().Fatalw(msg, keysAndValues...)
}

// Sync memory data to log files.
func (l *Logger) Sync() {
	_ = l.zapLogger.Sync()
}
