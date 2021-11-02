package log

import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"
	"go.uber.org/zap/zapcore"
)

const (
	consoleFormat = "console"
	textFormat    = "text"
	jsonFormat    = "json"
)

const (
	flagName              = "log.name"
	flagLevel             = "log.level"
	flagFormat            = "log.format"
	flagDisableColor      = "log.disable-color"
	flagDisableCaller     = "log.disable-caller"
	flagDisableStacktrace = "log.disable-stacktrace"
	flagOutputPaths       = "log.output-paths"
	flagErrorOutputPaths  = "log.error-output-paths"
	flagEnableRotate      = "log.enable-rotate"
	flagRotateMaxSize     = "log.rotate-max-size"
	flagRotateMaxAge      = "log.rotate-max-age"
	flagRotateMaxBackups  = "log.rotate-max-backups"
	flagRotateLocaltime   = "log.rotate-localtime"
	flagRotateCompress    = "log.rotate-compress"
)

var defaultOpts = NewOptions()

// Options for logger.
type Options struct {
	Name              string   `json:"name"               mapstructure:"name"`
	Level             string   `json:"level"              mapstructure:"level"`
	Format            string   `json:"format"             mapstructure:"format"`
	DisableColor      bool     `json:"disable-color"      mapstructure:"disable-color"`
	DisableCaller     bool     `json:"disable-caller"     mapstructure:"disable-caller"`
	DisableStacktrace bool     `json:"disable-stacktrace" mapstructure:"disable-stacktrace"`
	OutputPaths       []string `json:"output-paths"       mapstructure:"output-paths"`
	ErrorOutputPaths  []string `json:"error-output-paths" mapstructure:"error-output-paths"`
	EnableRotate      bool     `json:"enable-rotate"      mapstructure:"enable-rotate"`
	RotateMaxSize     int      `json:"rotate-max-size"    mapstructure:"rotate-max-size"`
	RotateMaxAge      int      `json:"rotate-max-age"     mapstructure:"rotate-max-age"`
	RotateMaxBackups  int      `json:"rotate-max-backups" mapstructure:"rotate-max-backups"`
	RotateLocaltime   bool     `json:"rotate-localtime"   mapstructure:"rotate-localtime"`
	RotateCompress    bool     `json:"rotate-compression" mapstructure:"rotate-compression"`
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

// AddFlagsTo a instace of pflag.FlagSet
func (o *Options) AddFlagsTo(fs *pflag.FlagSet) {
	fs.StringVar(&o.Name, flagName, o.Name, "The name of the logger")
	fs.StringVar(&o.Level, flagLevel, o.Level, "Log level")
	fs.StringVar(
		&o.Format,
		flagFormat,
		o.Format,
		"Log format, console/text or json",
	)
	fs.BoolVar(
		&o.DisableColor,
		flagDisableColor,
		o.DisableColor,
		"Disable ansi colors in console/text format logs",
	)
	fs.BoolVar(
		&o.DisableCaller,
		flagDisableCaller,
		o.DisableCaller,
		"Disable log caller that print line number of the log entry",
	)
	fs.BoolVar(
		&o.DisableStacktrace,
		flagDisableStacktrace,
		o.DisableStacktrace,
		"Disable log stack trace for logs at or above panic level",
	)
	fs.StringSliceVar(
		&o.OutputPaths,
		flagOutputPaths,
		o.OutputPaths,
		"Log files or stdout/stderr that contain all level entries",
	)
	fs.StringSliceVar(
		&o.ErrorOutputPaths,
		flagErrorOutputPaths,
		o.ErrorOutputPaths,
		"Log files or stdout/stderr that only contain logger internal errors",
	)
	fs.BoolVar(
		&o.EnableRotate,
		flagEnableRotate,
		o.EnableRotate,
		"Enable log rotation or not",
	)
	fs.IntVar(
		&o.RotateMaxSize,
		flagRotateMaxSize,
		o.RotateMaxSize,
		"The maximum size in megabytes of the log file before it gets rotated",
	)
	fs.IntVar(
		&o.RotateMaxAge,
		flagRotateMaxAge,
		o.RotateMaxAge,
		"The maximum number of days to retain old log files based on the timestamp encoded in their filename",
	)
	fs.IntVar(
		&o.RotateMaxBackups,
		flagRotateMaxBackups,
		o.RotateMaxBackups,
		"The maximum number of old log files to retain",
	)
	fs.BoolVar(
		&o.RotateLocaltime,
		flagRotateLocaltime,
		o.RotateLocaltime,
		"The timestamps in backup files name use the server local time instead of UTC time or not",
	)
	fs.BoolVar(
		&o.RotateCompress,
		flagRotateCompress,
		o.RotateCompress,
		"Using gzip to compress the rotated log files or not",
	)
}

// Validate the options fields.
func (o *Options) Validate() []error {
	var errs []error

	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(o.Level)); err != nil {
		errs = append(errs, err)
	}

	format := strings.ToLower(o.Format)
	if format != consoleFormat && format != textFormat && format != jsonFormat {
		errs = append(errs, fmt.Errorf("invalid log format: %q", o.Format))
	}

	return errs
}
