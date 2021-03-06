package app

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"gobackend/pkg/errors"
	cliflag "gobackend/pkg/flag"
	"gobackend/pkg/log"
	"gobackend/pkg/term"
	"gobackend/pkg/version"
	"gobackend/pkg/version/verflag"
)

// RunFunc defines the application's startup callback function.
type RunFunc func(binaryName string) error

// Option defines optional parameters for initializing the application structure.
type Option func(*App)

// WithOptions to open the application's function to read from the command line
// or read parameters from the configuration file.
func WithOptions(opt CliOptions) Option {
	return func(a *App) {
		a.options = opt
	}
}

// WithRunFunc is used to set the application startup callback function option.
func WithRunFunc(run RunFunc) Option {
	return func(a *App) {
		a.runFunc = run
	}
}

// WithDescription is used to set the description of the application.
func WithDescription(desc string) Option {
	return func(a *App) {
		a.description = desc
	}
}

// WithSilence sets the application to silent mode, in which the program startup
// information, configuration information, and version information are not
// printed in the console.
func WithSilence() Option {
	return func(a *App) {
		a.silence = true
	}
}

// WithProcessLock make sure only one process is running at a time.
func WithProcessLock(pidDir string) Option {
	return func(a *App) {
		a.processLock = true
		a.pidDir = pidDir
	}
}

// WithRunModeEnv custom your own environment variable to specify run environments(dev/test/prod).
func WithRunModeEnv(runModeEnv string) Option {
	return func(a *App) {
		a.runModeEnv = runModeEnv
	}
}

// WithNoVersion set the application does not provide version flag.
func WithNoVersion() Option {
	return func(a *App) {
		a.noVersion = true
	}
}

// WithNoConfig set the application does not provide config flag.
func WithNoConfig() Option {
	return func(a *App) {
		a.noConfig = true
	}
}

// WithValidArgs set the validation function to valid non-flag arguments.
func WithValidArgs(args cobra.PositionalArgs) Option {
	return func(a *App) {
		a.args = args
	}
}

// WithDefaultValidArgs set default validation function to valid non-flag arguments.
func WithDefaultValidArgs() Option {
	return func(a *App) {
		a.args = func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				if len(arg) > 0 {
					return fmt.Errorf("%q does not take any arguments, got %q", cmd.CommandPath(), args)
				}
			}

			return nil
		}
	}
}

// App is the main structure of a cli application.
// It is recommended that an app be created with the app.NewApp() function.
type App struct {
	name        string
	binaryName  string
	description string
	options     CliOptions
	runFunc     RunFunc
	runModeEnv  string
	processLock bool
	pidDir      string
	silence     bool
	noVersion   bool
	noConfig    bool
	commands    []*Command
	args        cobra.PositionalArgs
	cmd         *cobra.Command
}

// New creates a new application instance based on the given application name,
// binary name, and other options.
func New(name string, binaryName string, opts ...Option) *App {
	a := &App{
		name:       name,
		binaryName: binaryName,
		runModeEnv: "RUN_MODE",
	}

	for _, o := range opts {
		o(a)
	}

	a.buildCommand()

	return a
}

// Run is used to launch the application.
func (a *App) Run() {
	if a.processLock {
		lock, lockFile, err := processLock(a.pidDir)
		if err != nil {
			fmt.Printf("%v\n", err)

			return
		}

		defer os.Remove(lockFile)
		defer lock.Close()
	}

	if err := a.cmd.Execute(); err != nil {
		fmt.Printf("%v %v\n", color.RedString("Error:"), err)
		os.Exit(1)
	}
}

// Command returns cobra command instance inside the application.
func (a *App) Command() *cobra.Command {
	return a.cmd
}

func (a *App) buildCommand() {
	cmd := cobra.Command{
		Use:   FormatBinaryName(a.binaryName),
		Short: a.name,
		Long:  a.description,
		// Stop printing usage when the command errors or not.
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          a.args,
	}

	cmd.SetOut(os.Stdout)
	cmd.SetErr(os.Stderr)
	cmd.Flags().SortFlags = true

	cliflag.InitFlags(cmd.Flags())

	if len(a.commands) > 0 {
		for _, command := range a.commands {
			cmd.AddCommand(command.cobraCommand())
		}
		cmd.SetHelpCommand(helpCommand(a.name))
	}

	if a.runFunc != nil {
		cmd.RunE = a.runCommand
	}

	var namedFlagSets cliflag.NamedFlagSets
	if a.options != nil {
		namedFlagSets = a.options.Flags()
		fs := cmd.Flags()
		for _, f := range namedFlagSets.FlagSets {
			fs.AddFlagSet(f)
		}

		usageFmt := "Usage:\n  %s\n"
		cols, _, _ := term.TerminalSize(cmd.OutOrStdout())
		cmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
			fmt.Fprintf(cmd.OutOrStdout(), "%s\n\n"+usageFmt, cmd.Long, cmd.UseLine())
			cliflag.PrintSections(cmd.OutOrStdout(), namedFlagSets, cols)
		})
		cmd.SetUsageFunc(func(cmd *cobra.Command) error {
			fmt.Fprintf(cmd.OutOrStderr(), usageFmt, cmd.UseLine())
			cliflag.PrintSections(cmd.OutOrStderr(), namedFlagSets, cols)

			return nil
		})
	}

	if !a.noVersion {
		verflag.AddFlags(namedFlagSets.FlagSet("global"))
	}

	if !a.noConfig {
		addConfigFlag(namedFlagSets.FlagSet("global"))
	}

	cliflag.AddGlobalFlags(namedFlagSets.FlagSet("global"), cmd.Name())

	a.cmd = &cmd
}

func (a *App) runCommand(cmd *cobra.Command, args []string) error {
	cliflag.PrintFlags(cmd.Flags())

	if !a.noVersion {
		// display application version information
		verflag.PrintAndExitIfRequested()
	}

	parseConfigFile(a.binaryName, a.runModeEnv)

	if !a.noConfig {
		if err := viper.BindPFlags(cmd.Flags()); err != nil {
			return err
		}

		if err := viper.Unmarshal(a.options); err != nil {
			return err
		}
	}

	logOptions := &log.Options{
		Name:              viper.GetString("log.name"),
		Level:             viper.GetString("log.level"),
		Format:            viper.GetString("log.format"),
		DisableColor:      viper.GetBool("log.disable-color"),
		DisableCaller:     viper.GetBool("log.disable-caller"),
		DisableStacktrace: viper.GetBool("log.disable-stacktrace"),
		OutputPaths:       viper.GetStringSlice("log.output-paths"),
		ErrorOutputPaths:  viper.GetStringSlice("log.error-output-paths"),
		EnableRotate:      viper.GetBool("log.enable-rotate"),
		RotateMaxSize:     viper.GetInt("log.rotate-max-size"),
		RotateMaxAge:      viper.GetInt("log.rotate-max-age"),
		RotateMaxBackups:  viper.GetInt("log.rotate-max-backups"),
		RotateLocaltime:   viper.GetBool("log.rotate-localtime"),
		RotateCompress:    viper.GetBool("log.rotate-compress"),
	}

	log.Init(logOptions)
	defer log.Sync()

	if a.options != nil {
		if err := a.applyOptionRules(); err != nil {
			return err
		}
	}

	if !a.silence {
		printWorkingDir()

		log.Infof("Starting %s ...", a.name)

		if !a.noVersion {
			log.Infof("Version: %s", version.Get().ToJSON())
		}

		if !a.noConfig {
			if a.runModeEnv != "" {
				log.Infof(
					"Run environment variable: %s, value: %s",
					a.runModeEnv,
					os.Getenv(a.runModeEnv),
				)
			}

			log.Infof("Config file used: %s", viper.ConfigFileUsed())
		}
	}

	if a.runFunc != nil {
		return a.runFunc(a.binaryName)
	}

	return nil
}

func (a *App) applyOptionRules() error {
	if completeableOptions, ok := a.options.(CompleteableOptions); ok {
		if err := completeableOptions.Complete(); err != nil {
			return err
		}
	}

	if errs := a.options.Validate(); len(errs) != 0 {
		return errors.NewAggregate(errs)
	}

	if printableOptions, ok := a.options.(PrintableOptions); ok && !a.silence {
		log.Infof("Config contents: %s", printableOptions.String())
	}

	return nil
}

func printWorkingDir() {
	wd, _ := os.Getwd()
	log.Infof("Working dir: %s", wd)
}
