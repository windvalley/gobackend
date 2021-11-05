package log

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// The default global logger instance.
	std = New(defaultOpts)
	mu  sync.Mutex
)

// Init package default logger instance with given options.
func Init(opts *Options) {
	mu.Lock()
	defer mu.Unlock()
	std = New(opts)
	resetDefaultLogger()
}

// SetHooks and replace global logger.
//
// Usage e.g.:
//
// monitorHook := func(entry log.Entry) error {
//	  if entry.Level >= log.ErrorLevel {
//        fmt.Println("alert!")
//    }
//    return nil
// }
//
// log.SetHooks(monitorHook)
//
// log.Error("server failed")
//
func SetHooks(hooks ...Hook) {
	zapLogger := std.zapLogger.WithOptions(zap.Hooks(hooks...))
	std = newLogger(zapLogger, std.options)
	resetDefaultLogger()
}

// New logger instance with given options.
func New(opts *Options) *Logger {
	if opts == nil {
		opts = NewOptions()
	}

	format := defaultOpts.Format
	if opts.Format != "" {
		format = strings.ToLower(opts.Format)
	}

	if len(opts.OutputPaths) == 0 {
		opts.OutputPaths = defaultOpts.OutputPaths
	}

	if len(opts.ErrorOutputPaths) == 0 {
		opts.ErrorOutputPaths = defaultOpts.ErrorOutputPaths
	}

	var level Level
	if err := level.UnmarshalText([]byte(opts.Level)); err != nil {
		level = InfoLevel
	}

	var development bool
	if level < InfoLevel {
		development = true
	}

	encodeLevel := zapcore.CapitalLevelEncoder
	if format == consoleFormat && !opts.DisableColor {
		encodeLevel = zapcore.CapitalColorLevelEncoder
	}

	var syncers, errSyncers []zapcore.WriteSyncer

	if opts.EnableRotate {
		for _, path := range opts.OutputPaths {
			syncer := getRotateSyncer(path, opts)
			syncers = append(syncers, syncer)
		}

		for _, path := range opts.ErrorOutputPaths {
			syncer := getRotateSyncer(path, opts)
			errSyncers = append(errSyncers, syncer)
		}
	} else {
		for _, path := range opts.OutputPaths {
			syncer := getNormalSyncer(path)
			syncers = append(syncers, syncer)
		}

		for _, path := range opts.ErrorOutputPaths {
			syncer := getNormalSyncer(path)
			errSyncers = append(errSyncers, syncer)
		}
	}

	sink := zapcore.NewMultiWriteSyncer(syncers...)
	errSink := zapcore.NewMultiWriteSyncer(errSyncers...)

	encoderConfig := zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "timestamp",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    encodeLevel,
		EncodeTime:     timeEncoder,
		EncodeDuration: milliSecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	loggerConfig := &zap.Config{
		Level:             zap.NewAtomicLevelAt(level),
		Development:       development,
		DisableCaller:     opts.DisableCaller,
		DisableStacktrace: opts.DisableStacktrace,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		OutputPaths:      opts.OutputPaths,
		ErrorOutputPaths: opts.ErrorOutputPaths,
	}

	encoder := zapcore.NewJSONEncoder(encoderConfig)
	if format == consoleFormat {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	l := zap.New(
		zapcore.NewCore(encoder, sink, zap.NewAtomicLevelAt(level)),
		buildOptions(loggerConfig, errSink)...,
	)

	l = l.WithOptions(zap.AddStacktrace(zapcore.PanicLevel), zap.AddCallerSkip(1))

	return newLogger(l, opts)
}

// NewStdInfoLogger returns *log.Logger which writes to std.zapLogger at info level.
func NewStdInfoLogger() *log.Logger {
	if std == nil {
		return nil
	}

	l, err := zap.NewStdLogAt(std.zapLogger, zapcore.InfoLevel)
	if err != nil {
		return nil
	}

	return l
}

func buildOptions(cfg *zap.Config, errSink zapcore.WriteSyncer) []zap.Option {
	opts := []zap.Option{zap.ErrorOutput(errSink)}

	if cfg.Development {
		opts = append(opts, zap.Development())
	}

	if !cfg.DisableCaller {
		opts = append(opts, zap.AddCaller())
	}

	stackLevel := ErrorLevel
	if cfg.Development {
		stackLevel = WarnLevel
	}
	if !cfg.DisableStacktrace {
		opts = append(opts, zap.AddStacktrace(stackLevel))
	}

	if scfg := cfg.Sampling; scfg != nil {
		opts = append(opts, zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			var samplerOpts []zapcore.SamplerOption
			if scfg.Hook != nil {
				samplerOpts = append(samplerOpts, zapcore.SamplerHook(scfg.Hook))
			}
			return zapcore.NewSamplerWithOptions(
				core,
				time.Second,
				cfg.Sampling.Initial,
				cfg.Sampling.Thereafter,
				samplerOpts...,
			)
		}))
	}

	if len(cfg.InitialFields) > 0 {
		fs := make([]Field, 0, len(cfg.InitialFields))
		keys := make([]string, 0, len(cfg.InitialFields))
		for k := range cfg.InitialFields {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fs = append(fs, Any(k, cfg.InitialFields[k]))
		}
		opts = append(opts, zap.Fields(fs...))
	}

	return opts
}

func getNormalSyncer(path string) (syncer zapcore.WriteSyncer) {
	switch path {
	case "stdout":
		syncer = zapcore.AddSync(os.Stdout)
	case "stderr":
		syncer = zapcore.AddSync(os.Stderr)
	default:
		file, err := openFile(path)
		if err != nil {
			panic(err)
		}
		syncer = zapcore.AddSync(file)
	}

	return
}

func getRotateSyncer(path string, opts *Options) (syncer zapcore.WriteSyncer) {
	switch path {
	case "stdout":
		syncer = zapcore.AddSync(os.Stdout)
	case "stderr":
		syncer = zapcore.AddSync(os.Stderr)
	default:
		rl := &lumberjack.Logger{
			Filename:   path,
			MaxSize:    opts.RotateMaxSize,
			MaxAge:     opts.RotateMaxAge,
			MaxBackups: opts.RotateMaxBackups,
			LocalTime:  opts.RotateLocaltime,
			Compress:   opts.RotateCompress,
		}
		syncer = zapcore.AddSync(rl)
	}

	return
}

func openFile(filename string) (*os.File, error) {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		err := os.MkdirAll(filepath.Dir(filename), 0744)
		if err != nil {
			return nil, fmt.Errorf("make directory for new logfile failed: %s", err)
		}

		return os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	}
	if err != nil {
		return nil, err
	}

	return os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
}
