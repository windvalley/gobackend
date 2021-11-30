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

	"gobackend/pkg/util"
)

const configFlagName = "config"

var cfgFile string

//nolint: gochecknoinits
func init() {
	pflag.StringVarP(&cfgFile, "config", "c", cfgFile, "Read configuration from a specified file, "+
		"support formats: JSON, TOML, YAML")
}

// addConfigFlag for a specific server to the specified FlagSet object.
func addConfigFlag(fs *pflag.FlagSet) {
	fs.AddFlag(pflag.Lookup(configFlagName))
}

func parseConfigFile(binaryName, runModeEnv string) {
	viper.AutomaticEnv()
	viper.SetEnvPrefix(strings.Replace(strings.ToUpper(binaryName), "-", "_", -1))
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		setConfigPath(binaryName)
		setConfigName(binaryName, runModeEnv)
	}

	if err := viper.ReadInConfig(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s read config file failed: %s\n", color.RedString("Error:"), err)
		os.Exit(1)
	}
}

func setConfigPath(binaryName string) {
	viper.AddConfigPath(".")
	viper.AddConfigPath("./configs")

	if names := strings.Split(binaryName, "-"); len(names) > 1 {
		viper.AddConfigPath(filepath.Join(util.HomeDir(), "."+names[0]))
	} else {
		viper.AddConfigPath(filepath.Join(util.HomeDir(), "."+binaryName))
	}
}

func setConfigName(binaryName, runModeEnv string) {
	runEnv, ok := os.LookupEnv(runModeEnv)
	if !ok {
		_, _ = fmt.Fprintf(
			os.Stderr,
			"%s env %s not exist, it will load development config file by default\n",
			color.YellowString("Warning:"),
			color.BlueString(runModeEnv),
		)

		viper.SetConfigName("dev." + binaryName)
	} else {
		switch runEnv {
		case "test":
			viper.SetConfigName("test." + binaryName)
		case "prod":
			viper.SetConfigName("prod." + binaryName)
		case "dev":
			viper.SetConfigName("dev." + binaryName)
		default:
			_, _ = fmt.Fprintf(
				os.Stderr,
				"%s invalid %s value, it will load development config file by default\n",
				color.YellowString("Warning:"),
				color.BlueString(runModeEnv),
			)

			viper.SetConfigName("dev." + binaryName)
		}
	}
}

//nolint: deadcode,unused
func printConfig() {
	keys := viper.AllKeys()

	if len(keys) > 0 {
		fmt.Printf("Configuration items:\n")

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
