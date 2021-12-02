package options

import (
	"encoding/json"

	"gobackend/internal/pkg/server"
	cliflag "gobackend/pkg/flag"

	genericoptions "gobackend/internal/pkg/options"
)

// Options ...
type Options struct {
	GenericServerRun *genericoptions.ServerRunOptions       `json:"server"   mapstructure:"server"`
	InsecureServing  *genericoptions.InsecureServingOptions `json:"insecure" mapstructure:"insecure"`
	SecureServing    *genericoptions.SecureServingOptions   `json:"secure"   mapstructure:"secure"`
	MySQL            *genericoptions.MySQLOptions           `json:"mysql"    mapstructure:"mysql"`
	Feature          *genericoptions.FeatureOptions         `json:"feature"  mapstructure:"feature"`
	Log              *genericoptions.LogOptions             `json:"log"      mapstructure:"log"`
}

// New creates a new Options object with default parameters.
func New() *Options {
	o := Options{
		GenericServerRun: genericoptions.NewServerRunOptions(),
		InsecureServing:  genericoptions.NewInsecureServingOptions(),
		SecureServing:    genericoptions.NewSecureServingOptions(),
		MySQL:            genericoptions.NewMySQLOptions(),
		Feature:          genericoptions.NewFeatureOptions(),
		Log:              genericoptions.NewLogOptions(),
	}

	return &o
}

// ApplyTo applies the run options to the method receiver and returns self.
func (o *Options) ApplyTo(c *server.Config) (lastErr error) {
	if lastErr = o.GenericServerRun.ApplyTo(c); lastErr != nil {
		return
	}

	if lastErr = o.InsecureServing.ApplyTo(c); lastErr != nil {
		return
	}

	if lastErr = o.SecureServing.ApplyTo(c); lastErr != nil {
		return
	}

	if lastErr = o.Feature.ApplyTo(c); lastErr != nil {
		return
	}

	return nil
}

// Flags returns flags for a specific APIServer by section name.
func (o *Options) Flags() (fss cliflag.NamedFlagSets) {
	o.GenericServerRun.AddFlags(fss.FlagSet("generic"))
	o.InsecureServing.AddFlags(fss.FlagSet("insecure serving"))
	o.SecureServing.AddFlags(fss.FlagSet("secure serving"))
	o.MySQL.AddFlags(fss.FlagSet("mysql"))
	o.Feature.AddFlags(fss.FlagSet("features"))
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

	errs = append(errs, o.GenericServerRun.Validate()...)
	errs = append(errs, o.InsecureServing.Validate()...)
	errs = append(errs, o.SecureServing.Validate()...)
	errs = append(errs, o.MySQL.Validate()...)
	errs = append(errs, o.Feature.Validate()...)
	errs = append(errs, o.Log.Validate()...)

	return errs
}
