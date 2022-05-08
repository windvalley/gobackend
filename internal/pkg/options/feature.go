package options

import (
	"github.com/spf13/pflag"

	"gobackend/internal/pkg/server"
)

// FeatureOptions contains configuration items related to API server features.
type FeatureOptions struct {
	EnableProfiling        bool `json:"profiling"      mapstructure:"profiling"`
	EnableMetrics          bool `json:"enable-metrics" mapstructure:"enable-metrics"`
	EnableOperationLogging bool `json:"operation-logging" mapstructure:"operation-logging"`
}

// NewFeatureOptions creates a FeatureOptions object with default parameters.
func NewFeatureOptions() *FeatureOptions {
	defaults := server.NewConfig()

	return &FeatureOptions{
		EnableMetrics:          defaults.EnableMetrics,
		EnableProfiling:        defaults.EnableProfiling,
		EnableOperationLogging: defaults.EnableOperationLogging,
	}
}

// ApplyTo applies the run options to the method receiver and returns self.
func (o *FeatureOptions) ApplyTo(c *server.Config) error {
	c.EnableProfiling = o.EnableProfiling
	c.EnableMetrics = o.EnableMetrics
	c.EnableOperationLogging = o.EnableOperationLogging

	return nil
}

// Validate is used to parse and validate the parameters entered by the user at
// the command line when the program starts.
func (o *FeatureOptions) Validate() []error {
	return nil
}

// AddFlags adds flags related to features for a specific api server to the
// specified FlagSet.
func (o *FeatureOptions) AddFlags(fs *pflag.FlagSet) {
	if fs == nil {
		return
	}

	fs.BoolVar(&o.EnableProfiling, "feature.profiling", o.EnableProfiling,
		"Enable profiling via web interface host:port/debug/pprof/")

	fs.BoolVar(&o.EnableMetrics, "feature.enable-metrics", o.EnableMetrics,
		"Enables metrics on the apiserver at /metrics")

	fs.BoolVar(
		&o.EnableOperationLogging,
		"feature.operation-logging",
		o.EnableOperationLogging,
		"Enable operation logging of the apiserver",
	)
}
