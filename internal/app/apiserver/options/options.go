package options

import (
	"encoding/json"

	genericoptions "gobackend/internal/pkg/options"
	"gobackend/internal/pkg/server"
	cliflag "gobackend/pkg/flag"
	"gobackend/pkg/log"
)

// Options ...
type Options struct {
	GenericServerRunOptions *genericoptions.ServerRunOptions       `json:"server"   mapstructure:"server"`
	InsecureServing         *genericoptions.InsecureServingOptions `json:"insecure" mapstructure:"insecure"`
	SecureServing           *genericoptions.SecureServingOptions   `json:"secure"   mapstructure:"secure"`
	MySQLOptions            *genericoptions.MySQLOptions           `json:"mysql"    mapstructure:"mysql"`
	FeatureOptions          *genericoptions.FeatureOptions         `json:"feature"  mapstructure:"feature"`
	Log                     *log.Options                           `json:"log"      mapstructure:"log"`
}

// New creates a new Options object with default parameters.
func New() *Options {
	o := Options{
		GenericServerRunOptions: genericoptions.NewServerRunOptions(),
		InsecureServing:         genericoptions.NewInsecureServingOptions(),
		SecureServing:           genericoptions.NewSecureServingOptions(),
		MySQLOptions:            genericoptions.NewMySQLOptions(),
		FeatureOptions:          genericoptions.NewFeatureOptions(),
		Log:                     log.NewOptions(),
	}

	return &o
}

// ApplyTo applies the run options to the method receiver and returns self.
func (o *Options) ApplyTo(c *server.Config) (lastErr error) {
	if lastErr = o.GenericServerRunOptions.ApplyTo(c); lastErr != nil {
		return
	}

	if lastErr = o.InsecureServing.ApplyTo(c); lastErr != nil {
		return
	}

	if lastErr = o.SecureServing.ApplyTo(c); lastErr != nil {
		return
	}

	if lastErr = o.FeatureOptions.ApplyTo(c); lastErr != nil {
		return
	}

	return nil
}

// Flags returns flags for a specific APIServer by section name.
func (o *Options) Flags() (fss cliflag.NamedFlagSets) {
	o.GenericServerRunOptions.AddFlags(fss.FlagSet("generic"))
	o.InsecureServing.AddFlags(fss.FlagSet("insecure serving"))
	o.SecureServing.AddFlags(fss.FlagSet("secure serving"))
	o.MySQLOptions.AddFlags(fss.FlagSet("mysql"))
	o.FeatureOptions.AddFlags(fss.FlagSet("features"))
	o.Log.AddFlagsTo(fss.FlagSet("logs"))

	return fss
}

func (o *Options) String() string {
	data, _ := json.Marshal(o)

	return string(data)
}

// Complete set default Options.
func (o *Options) Complete() error {
	return nil
}

// Validate checks Options and return a slice of found errs.
func (o *Options) Validate() []error {
	var errs []error

	errs = append(errs, o.GenericServerRunOptions.Validate()...)
	errs = append(errs, o.InsecureServing.Validate()...)
	errs = append(errs, o.SecureServing.Validate()...)
	errs = append(errs, o.MySQLOptions.Validate()...)
	errs = append(errs, o.FeatureOptions.Validate()...)
	errs = append(errs, o.Log.Validate()...)

	return errs
}
