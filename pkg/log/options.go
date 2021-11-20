package log

import "go.uber.org/zap/zapcore"

const consoleFormat = "console"

var defaultOpts = NewOptions()

// Options for logger.
type Options struct {
	Name              string
	Level             string
	Format            string
	DisableColor      bool
	DisableCaller     bool
	DisableStacktrace bool
	OutputPaths       []string
	ErrorOutputPaths  []string
	EnableRotate      bool
	RotateMaxSize     int
	RotateMaxAge      int
	RotateMaxBackups  int
	RotateLocaltime   bool
	RotateCompress    bool
}

// NewOptions return default log options instance.
func NewOptions() *Options {
	return &Options{
		Name:              "",
		Level:             zapcore.InfoLevel.String(),
		Format:            consoleFormat,
		DisableColor:      false,
		DisableCaller:     false,
		DisableStacktrace: false,
		OutputPaths: []string{
			"stdout",
		},
		ErrorOutputPaths: []string{
			"stderr",
		},
		EnableRotate:     false,
		RotateMaxSize:    100,
		RotateMaxAge:     28,
		RotateMaxBackups: 0,
		RotateLocaltime:  true,
		RotateCompress:   false,
	}
}
