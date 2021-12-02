package options

import (
	"fmt"

	"github.com/spf13/pflag"

	"gobackend/internal/pkg/middleware"
	"gobackend/internal/pkg/server"
)

// ServerRunOptions contains the options while running a generic api server.
type ServerRunOptions struct {
	Mode        string   `json:"mode"        mapstructure:"mode"`
	Healthz     bool     `json:"healthz"     mapstructure:"healthz"`
	Middlewares []string `json:"middlewares" mapstructure:"middlewares"`
}

// NewServerRunOptions creates a new ServerRunOptions object with default parameters.
func NewServerRunOptions() *ServerRunOptions {
	defaults := server.NewConfig()

	return &ServerRunOptions{
		Mode:        defaults.Mode,
		Healthz:     defaults.Healthz,
		Middlewares: defaults.Middlewares,
	}
}

// ApplyTo applies the run options to the method receiver and returns self.
func (s *ServerRunOptions) ApplyTo(c *server.Config) error {
	c.Mode = s.Mode
	c.Healthz = s.Healthz
	c.Middlewares = s.Middlewares

	return nil
}

// Validate checks validation of ServerRunOptions.
func (s *ServerRunOptions) Validate() (errs []error) {
	if s.Mode != "debug" && s.Mode != "test" && s.Mode != "release" {
		errs = append(errs, fmt.Errorf(
			"unknown server.mode: %s, available mode: [debug release test]",
			s.Mode,
		))
	}

	var availableMiddlewares, invalidMiddlewares []string

	for _, m := range s.Middlewares {
		if _, ok := middleware.Middlewares[m]; !ok {
			invalidMiddlewares = append(invalidMiddlewares, m)
		}
	}

	if len(invalidMiddlewares) != 0 {
		for m := range middleware.Middlewares {
			availableMiddlewares = append(availableMiddlewares, m)
		}

		errs = append(errs, fmt.Errorf(
			"unknown server.middlewares: %v, available middlewares: %v",
			invalidMiddlewares,
			availableMiddlewares,
		))
	}

	return
}

// AddFlags adds flags for a specific APIServer to the specified FlagSet.
func (s *ServerRunOptions) AddFlags(fs *pflag.FlagSet) {
	// Note: the weird ""+ in below lines seems to be the only way to get gofmt to
	// arrange these text blocks sensibly. Grrr.
	fs.StringVar(&s.Mode, "server.mode", s.Mode, ""+
		"Start the server in a specified server mode. Supported server mode: debug, test, release.")

	fs.BoolVar(&s.Healthz, "server.healthz", s.Healthz, ""+
		"Add self readiness check and install /healthz router.")

	fs.StringSliceVar(&s.Middlewares, "server.middlewares", s.Middlewares, ""+
		"List of allowed middlewares for server, comma separated. If this list is empty default middlewares will be used.")
}
