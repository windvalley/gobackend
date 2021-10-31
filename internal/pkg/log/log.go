package log

import (
	"github.com/go-project-pkg/log"

	"go-web-demo/internal/pkg/options"
)

// Init ...
func Init(opts *options.Log) {
	logOpts := &log.Options{
		Name:              opts.Name,   // logger name
		Level:             opts.Level,  // debug, info, warn, error, panic, dpanic, fatal
		Format:            opts.Format, // json, console/text
		DisableColor:      opts.DisableColor,
		DisableCaller:     opts.DisableCaller,
		DisableStacktrace: opts.DisableStacktrace,
		// Aplication's all levels logs.
		OutputPaths: opts.OutputPaths,
		// Only include zap internal errors, not include application's any level logs.
		ErrorOutputPaths: opts.ErrorOutputPaths,
		// Enable log files rotation feature or not.
		EnableRotate: true,
		// Take effect when EnableRotate is true.
		RotateOptions: &log.RotateOptions{
			// Maximum size in megabytes of the log file before it gets rotated.
			// Default: 100, if the value is 0, the log files will not be rotated.
			MaxSize: 1,
			// Saved days, default 0, means no limit.
			MaxAge: 30,
			// Saved count, default 0, means no limit.
			MaxBackups: 2,
			// Use local time in log file name, default false.
			LocalTime: true,
			// Gzip log files, default false.
			Compress: false,
		},
	}

	log.Init(logOpts)
}
