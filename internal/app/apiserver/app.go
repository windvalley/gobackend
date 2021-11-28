package apiserver

import (
	"gobackend/pkg/app"

	"gobackend/internal/app/apiserver/config"
	"gobackend/internal/app/apiserver/options"
)

const commandDesc = `Go Web Demo
`

// NewApp creates a App object with default parameters.
func NewApp(basename string) *app.App {
	opts := options.New()

	application := app.New("APIServer",
		basename,
		app.WithOptions(opts),
		app.WithDescription(commandDesc),
		app.WithProcessLock("/tmp"),
		app.WithDefaultValidArgs(),
		app.WithRunFunc(run(opts)),
	)

	return application
}

func run(opts *options.Options) app.RunFunc {
	return func(basename string) error {
		cfg := config.New(opts)

		server, err := createAPIServer(cfg)
		if err != nil {
			return err
		}

		return server.PrepareRun().Run()
	}
}
