package options

import (
	"fmt"

	"github.com/spf13/pflag"

	"gobackend/internal/pkg/server"
)

// SecureServingOptions contains configuration items related to HTTPS server startup.
type SecureServingOptions struct {
	BindAddress string `json:"bind-address" mapstructure:"bind-address"`
	// BindPort is ignored when Listener is set, will serve HTTPS even with 0.
	BindPort int `json:"bind-port"    mapstructure:"bind-port"`
	// Required set to true means that BindPort cannot be zero.
	Required bool
	// TLS cert info for serving secure traffic
	TLS TLS `json:"tls"          mapstructure:"tls"`
}

// TLS contains configuration items related to certificate.
type TLS struct {
	// CertFile is a file containing a PEM-encoded certificate, and possibly the complete certificate chain
	CertFile string `json:"cert-file"        mapstructure:"cert-file"`
	// KeyFile is a file containing a PEM-encoded private key for the certificate specified by CertFile
	KeyFile string `json:"key-file" mapstructure:"key-file"`
}

// NewSecureServingOptions creates a SecureServingOptions object with default parameters.
func NewSecureServingOptions() *SecureServingOptions {
	return &SecureServingOptions{
		BindAddress: "0.0.0.0",
		BindPort:    8443,
		Required:    false,
		TLS: TLS{
			CertFile: "",
			KeyFile:  "",
		},
	}
}

// ApplyTo generic server config
func (s *SecureServingOptions) ApplyTo(c *server.Config) error {
	c.SecureServing = &server.SecureServingInfo{
		BindAddress: s.BindAddress,
		BindPort:    s.BindPort,
		TLS: server.TLS{
			CertFile: s.TLS.CertFile,
			KeyFile:  s.TLS.KeyFile,
		},
	}

	return nil
}

// Validate is used to parse and validate the parameters entered by the user at
// the command line when the program starts.
func (s *SecureServingOptions) Validate() []error {
	if s == nil {
		return nil
	}

	errors := []error{}

	if s.Required && s.BindPort < 1 || s.BindPort > 65535 {
		errors = append(
			errors,
			fmt.Errorf(
				"--secure.bind-port %v must be between 1 and 65535, inclusive. It cannot be turned off with 0",
				s.BindPort,
			),
		)
	} else if s.BindPort < 0 || s.BindPort > 65535 {
		errors = append(
			errors,
			fmt.Errorf(
				"--secure.bind-port %v must be between 0 and 65535, inclusive. 0 for turning off secure port",
				s.BindPort,
			),
		)
	}

	return errors
}

// AddFlags adds flags related to HTTPS server for a specific APIServer to the
// specified FlagSet.
func (s *SecureServingOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&s.BindAddress, "secure.bind-address", s.BindAddress, ""+
		"The IP address on which to listen for the --secure.bind-port port. The "+
		"associated interface(s) must be reachable by the rest of the engine, and by CLI/web "+
		"clients. If blank, all interfaces will be used (0.0.0.0 for all IPv4 interfaces and :: for all IPv6 interfaces).")

	desc := "The port on which to serve HTTPS."

	if !s.Required {
		desc += " If 0, app will not enable HTTPS."
	}

	fs.IntVar(&s.BindPort, "secure.bind-port", s.BindPort, desc)

	fs.StringVar(&s.TLS.CertFile, "secure.tls.cert-file", s.TLS.CertFile, ""+
		"File containing the default x509 Certificate for HTTPS. (CA cert, if any, concatenated "+
		"after server cert).")

	fs.StringVar(&s.TLS.KeyFile,
		"secure.tls.key-file",
		s.TLS.KeyFile,
		"File containing the default x509 private key matching --secure.tls.cert-file.")
}
