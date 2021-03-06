package options

import (
	"fmt"
	"time"

	"github.com/spf13/pflag"

	"gobackend/internal/pkg/gormlog"
)

// MySQLOptions defines options for mysql database.
type MySQLOptions struct {
	Host                  string        `json:"host,omitempty"                     mapstructure:"host"`
	Username              string        `json:"username,omitempty"                 mapstructure:"username"`
	Password              string        `json:"-"                                  mapstructure:"password"`
	Database              string        `json:"database"                           mapstructure:"database"`
	MaxIdleConnections    int           `json:"max-idle-connections,omitempty"     mapstructure:"max-idle-connections"`
	MaxOpenConnections    int           `json:"max-open-connections,omitempty"     mapstructure:"max-open-connections"`
	MaxConnectionLifetime time.Duration `json:"max-connection-lifetime,omitempty"  mapstructure:"max-connection-lifetime"`
	LogLevel              string        `json:"log-level"                          mapstructure:"log-level"`
	AutoMigrate           bool          `json:"auto-migrate"                       mapstructure:"auto-migrate"`
}

// NewMySQLOptions create a `zero` value instance.
func NewMySQLOptions() *MySQLOptions {
	return &MySQLOptions{
		Host:                  "127.0.0.1:3306",
		Username:              "",
		Password:              "",
		Database:              "",
		MaxIdleConnections:    100,
		MaxOpenConnections:    100,
		MaxConnectionLifetime: time.Duration(120) * time.Second,
		LogLevel:              "silent",
		AutoMigrate:           false,
	}
}

// Validate verifies flags passed to MySQLOptions.
func (o *MySQLOptions) Validate() (errs []error) {
	if _, ok := gormlog.LogLevelMap[o.LogLevel]; !ok {
		errs = append(errs, fmt.Errorf(
			"unknown mysql.log-level: %s, available log level: [silent error warn info]",
			o.LogLevel,
		))
	}

	return
}

// AddFlags adds flags related to mysql storage for a specific APIServer to the specified FlagSet.
func (o *MySQLOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(
		&o.Host,
		"mysql.host",
		o.Host,
		"MySQL service host address. If left blank, the following related mysql options will be ignored.",
	)

	fs.StringVar(
		&o.Username,
		"mysql.username",
		o.Username,
		"Username for access to mysql service.",
	)

	fs.StringVar(
		&o.Password,
		"mysql.password",
		o.Password,
		"Password for access to mysql, should be used pair with password.",
	)

	fs.StringVar(
		&o.Database,
		"mysql.database",
		o.Database,
		"Database name for the server to use.",
	)

	fs.IntVar(
		&o.MaxIdleConnections,
		"mysql.max-idle-connections",
		o.MaxOpenConnections,
		"Maximum idle connections allowed to connect to mysql.",
	)

	fs.IntVar(
		&o.MaxOpenConnections,
		"mysql.max-open-connections",
		o.MaxOpenConnections,
		"Maximum open connections allowed to connect to mysql.",
	)

	fs.DurationVar(
		&o.MaxConnectionLifetime,
		"mysql.max-connection-lifetime",
		o.MaxConnectionLifetime,
		"Maximum connection life time allowed to connect to mysql.",
	)

	fs.StringVar(
		&o.LogLevel,
		"mysql.log-level",
		o.LogLevel,
		"Specify gorm log level: silent/error/warn/info.",
	)

	fs.BoolVar(
		&o.AutoMigrate,
		"mysql.auto-migrate",
		o.AutoMigrate,
		"Auto migrate database or not.",
	)
}
