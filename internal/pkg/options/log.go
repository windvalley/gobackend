package options

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/pflag"
	"go.uber.org/zap/zapcore"

	"gobackend/pkg/util"
)

const (
	consoleFormat = "console"
	jsonFormat    = "json"
)

const (
	flagLogName              = "log.name"
	flagLogLevel             = "log.level"
	flagLogFormat            = "log.format"
	flagLogDisableColor      = "log.disable-color"
	flagLogDisableCaller     = "log.disable-caller"
	flagLogDisableStacktrace = "log.disable-stacktrace"
	flagLogOutputPaths       = "log.output-paths"
	flagLogErrorOutputPaths  = "log.error-output-paths"
	flagLogEnableRotate      = "log.enable-rotate"
	flagLogRotateMaxSize     = "log.rotate-max-size"
	flagLogRotateMaxAge      = "log.rotate-max-age"
	flagLogRotateMaxBackups  = "log.rotate-max-backups"
	flagLogRotateLocaltime   = "log.rotate-localtime"
	flagLogRotateCompress    = "log.rotate-compress"
)

// LogOptions options.
type LogOptions struct {
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
	RotateCompress    bool     `json:"rotate-compress"    mapstructure:"rotate-compress"`
}

// NewLogOptions return default log options instance.
func NewLogOptions() *LogOptions {
	return &LogOptions{
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
func (o *LogOptions) AddFlagsTo(fs *pflag.FlagSet) {
	fs.StringVar(&o.Name, flagLogName, o.Name, "The name of the logger")
	fs.StringVar(&o.Level, flagLogLevel, o.Level, "Log level")
	fs.StringVar(
		&o.Format,
		flagLogFormat,
		o.Format,
		"Log format, console/text or json",
	)
	fs.BoolVar(
		&o.DisableColor,
		flagLogDisableColor,
		o.DisableColor,
		"Disable ansi colors in console/text format logs",
	)
	fs.BoolVar(
		&o.DisableCaller,
		flagLogDisableCaller,
		o.DisableCaller,
		"Disable log caller that print line number of the log entry",
	)
	fs.BoolVar(
		&o.DisableStacktrace,
		flagLogDisableStacktrace,
		o.DisableStacktrace,
		"Disable log stack trace for logs at or above panic level",
	)
	fs.StringSliceVar(
		&o.OutputPaths,
		flagLogOutputPaths,
		o.OutputPaths,
		"Log files or stdout/stderr that contain all level entries",
	)
	fs.StringSliceVar(
		&o.ErrorOutputPaths,
		flagLogErrorOutputPaths,
		o.ErrorOutputPaths,
		"Log files or stdout/stderr that only contain logger internal errors",
	)
	fs.BoolVar(
		&o.EnableRotate,
		flagLogEnableRotate,
		o.EnableRotate,
		"Enable log rotation or not",
	)
	fs.IntVar(
		&o.RotateMaxSize,
		flagLogRotateMaxSize,
		o.RotateMaxSize,
		"The maximum size in megabytes of the log file before it gets rotated",
	)
	fs.IntVar(
		&o.RotateMaxAge,
		flagLogRotateMaxAge,
		o.RotateMaxAge,
		"The maximum number of days to retain old log files based on the timestamp encoded in their filename",
	)
	fs.IntVar(
		&o.RotateMaxBackups,
		flagLogRotateMaxBackups,
		o.RotateMaxBackups,
		"The maximum number of old log files to retain",
	)
	fs.BoolVar(
		&o.RotateLocaltime,
		flagLogRotateLocaltime,
		o.RotateLocaltime,
		"The timestamps in backup files name use the server local time instead of UTC time or not",
	)
	fs.BoolVar(
		&o.RotateCompress,
		flagLogRotateCompress,
		o.RotateCompress,
		"Using gzip to compress the rotated log files or not",
	)
}

// Validate the options fields.
func (o *LogOptions) Validate() []error {
	var (
		errs            []error
		availableLevels = []string{"debug", "info", "warn", "error", "dpanic", "panic", "fatal"}
	)

	if !util.HasEntry(availableLevels, o.Level) {
		errs = append(errs, fmt.Errorf(
			"unknown log.level: %s, available log levels: %v",
			o.Level,
			availableLevels,
		))
	}

	format := strings.ToLower(o.Format)
	if format != consoleFormat && format != jsonFormat {
		errs = append(errs, fmt.Errorf(
			"unknown log.format: %s, available log formats: [%s %s]",
			o.Format,
			consoleFormat,
			jsonFormat,
		))
	}

	errs = checkOutputPaths(flagLogOutputPaths, o.OutputPaths, errs)
	errs = checkOutputPaths(flagLogErrorOutputPaths, o.ErrorOutputPaths, errs)

	return errs
}

func checkOutputPaths(flagName string, paths []string, errs []error) []error {
	for _, f := range paths {
		if strings.ToLower(f) == "stdout" || strings.ToLower(f) == "stderr" {
			if f != "stdout" && f != "stderr" {
				errs = append(errs, fmt.Errorf(
					"invalid %s: %s, should be: %s",
					flagName,
					f,
					strings.ToLower(f),
				))
			}

			continue
		}

		if err := checkLogfile(f); err != nil {
			errs = append(
				errs,
				fmt.Errorf("unavailable %s: %s, error: %w", flagName, f, err),
			)
		}
	}

	return errs
}

func checkLogfile(filename string) error {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		err = os.MkdirAll(filepath.Dir(filename), 0744)
		if err != nil {
			return fmt.Errorf("make directory for new logfile failed: %w", err)
		}

		_, err = os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)

		return err
	}

	if err != nil {
		return err
	}

	_, err = os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)

	return err
}
