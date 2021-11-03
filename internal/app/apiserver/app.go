package apiserver

import (
	"go-web-demo/pkg/app"
	"go-web-demo/pkg/log"

	"go-web-demo/internal/app/apiserver/config"
	"go-web-demo/internal/app/apiserver/options"
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
		app.WithDefaultValidArgs(),
		app.WithRunFunc(run(opts)),
	)

	return application
}

func run(opts *options.Options) app.RunFunc {
	return func(basename string) error {
		log.Init(opts.Log)
		defer log.Sync()

		cfg := config.New(opts)
		return Run(cfg)
	}
}

// Run runs the specified APIServer. This should never exit.
func Run(cfg *config.Config) error {
	server, err := createAPIServer(cfg)
	if err != nil {
		return err
	}

	return server.PrepareRun().Run()
}
