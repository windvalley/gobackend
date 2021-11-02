# go-project-pkg/log

Wrap [zap](https://github.com/uber-go/zap) for easy using.

## Installation

```sh
$ go get -u github.com/go-project-pkg/log
```

## Usage

Use default logger:

```go
import "github.com/go-project-pkg/log"

func main() {
    defer log.Sync()

    log.Info("Hello world!")
    log.Info("Hello ", log.String("string_key", "value"), log.Int("int_key", 666))
    log.Infof("Hello %s!", "world")
    log.Infow("Hello ", "string_key", "value", "int_key", 666)

    log.WithName("logger1").Warn("I am logger1")
    log.WithName("logger2").Warn("I am logger2")

    log.WithFields(log.String("f1", "value"), log.Int("f2", 888)).Error("Hello world!")
    log.WithName("logger3").WithFields(log.String("f1", "value"), log.Int("f2", 888)).Error("Hello world!")

    ctx := log.WithFields(String("f1", "value"), Int("f2", 888)).ToContext(context.Background())
    log.FromContext(ctx).Info("hello world!")
}
```

Custom your own logger:

```go
import "github.com/go-project-pkg/log"

func init() {
    opts := &log.Options{
        Name:              "",        // logger name
        Level:             "debug",   // debug, info, warn, error, panic, dpanic, fatal
        Format:            "console", // json, console/text
        DisableColor:      false,
        DisableCaller:     false,
        DisableStacktrace: false,
        // Aplication's all levels logs.
        OutputPaths: []string{
            "stdout", // os.Stdout
            "/var/log/app/app.log",
        },
        // Only include zap internal errors, not include application's any level logs.
        ErrorOutputPaths: []string{
            "stderr", // os.Stderr
            "/var/log/app/error.log",
        },
        // Enable log files rotation feature or not.
        EnableRotate: true,
        // Take effect when EnableRotate is true.
        RotateOptions: &log.RotateOptions{
            // Maximum size in megabytes of the log file before it gets rotated.
            // Default: 100, if the value is 0, the log files will not be rotated.
            MaxSize:    1,
            // Saved days, default 0, means no limit.
            MaxAge:     30,
            // Saved count, default 0, means no limit.
            MaxBackups: 2,
            // Use local time in log file name, default false.
            LocalTime:  true,
            // Gzip log files, default false.
            Compress:   false,
        },
    }

    log.Init(opts)
}

func main() {
    defer log.Sync()

    log.Info("Hello world!")
    log.Info("Hello ", log.String("string_key", "value"), log.Int("int_key", 666))
    log.Infof("Hello %s!", "world")
    log.Infow("Hello ", "string_key", "value", "int_key", 666)

    log.WithName("logger1").Warn("I am logger1")
    log.WithName("logger2").Warn("I am logger2")

    log.WithFields(log.String("f1", "value"), log.Int("f2", 888)).Error("Hello world!")
    log.WithName("logger3").WithFields(log.String("f1", "value"), log.Int("f2", 888)).Error("Hello world!")

    ctx := log.WithFields(String("f1", "value"), Int("f2", 888)).ToContext(context.Background())
    log.FromContext(ctx).Info("hello world!")

    // log files rotation test
    for i := 0; i <= 20000; i++ {
        log.Infof("hello world: %d", i)
    }
}
```

Use `log.C(ctx context.Context)` for getting logger with additional log fields by cooperating with gin's middleware:

```go
import "github.com/go-project-pkg/log"

// A middleware of gin for setting logger that with custom fileds to gin.Context
func Context() gin.HandlerFunc {
    return func(c *gin.Context) {
        l := log.WithFields(
            log.String("x-request-id", c.GetString(XRequestIDKey)),
            log.String("username", c.GetString(UsernameKey)),
        )
        c.Set(log.ContextLoggerName, l)

        c.Next()
    }
}

// Others place that use the logger.
func (u *UserController) Get(c *gin.Context) {
    // Get logger that with fileds from gin.Context and log a message.
    log.C(c).Debug("user get called")
}
```

You can add hooks to realize some useful features, like alerting when encountering error logs.

Use `log.SetHooks(hooks ...log.Hook)` for global logger:

```go
func main() {
    defer log.Sync()

    monitorHook1 := func(entry log.Entry) error {
        if entry.Level >= log.ErrorLevel {
            fmt.Println("hook1 alert!")
        }

        // This error is zap internal error, and it will write to 'ErrorOutputPaths'.
        return errors.New("alert hook failed")
    }

    monitorHook2 := func(entry log.Entry) error {
        if entry.Level >= log.ErrorLevel {
            fmt.Println("hook2 alert!")
        }

        return nil
    }

    log.SetHooks(monitorHook1, monitorHook2)

    log.Error("set hooks: server error")
}
```

Use `log.WithHooks(hooks ...log.Hook)` for current logger instance:

```go
func main() {
    defer log.Sync()

    monitorHook1 := func(entry log.Entry) error {
        if entry.Level >= log.ErrorLevel {
            fmt.Println("hook1 alert!")
        }

        // This error is zap internal error, and it will write to 'ErrorOutputPaths'.
        return errors.New("alert hook failed")
    }

    log.WithHooks(monitorHook1).Error("with hooks: server error")
}
```

## License

This project is under the MIT License. See the [LICENSE](LICENSE) file for the full license text.
