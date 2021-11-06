package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/gosuri/uitable"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"go-web-backend/pkg/util"
)

const configFlagName = "config"

var cfgFile string

//nolint: gochecknoinits
func init() {
	pflag.StringVarP(&cfgFile, "config", "c", cfgFile, "Read configuration from a specified file, "+
		"support formats: JSON, TOML, YAML, HCL")
}

// addConfigFlag for a specific server to the specified FlagSet object.
func addConfigFlag(fs *pflag.FlagSet) {
	fs.AddFlag(pflag.Lookup(configFlagName))
}

func parseConfigFile(basename string) {
	viper.AutomaticEnv()
	viper.SetEnvPrefix(strings.Replace(strings.ToUpper(basename), "-", "_", -1))
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	//cobra.OnInitialize(func() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")

		if names := strings.Split(basename, "-"); len(names) > 1 {
			viper.AddConfigPath(filepath.Join(util.HomeDir(), "."+names[0]))
		}

		viper.SetConfigName(basename)
	}

	if err := viper.ReadInConfig(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s read config file failed: %s\n", color.RedString("Error:"), err)
		os.Exit(1)
	}
	//})
}

//nolint: deadcode,unused
func printConfig() {
	keys := viper.AllKeys()
	if len(keys) > 0 {
		fmt.Println("Configuration items:")
		table := uitable.New()
		table.Separator = " "
		table.MaxColWidth = 80
		table.RightAlign(0)
		for _, k := range keys {
			table.AddRow(fmt.Sprintf("%s:", k), viper.Get(k))
		}
		fmt.Printf("%v", table)
	}
}
