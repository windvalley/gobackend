package config

import "gobackend/internal/app/apiserver/options"

// Config of apiserver.
type Config struct {
	*options.Options
}

// New config instance.
func New(opts *options.Options) *Config {
	return &Config{opts}
}
