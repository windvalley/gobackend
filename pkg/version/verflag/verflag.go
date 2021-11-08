package verflag

import (
	"fmt"
	"os"
	"strconv"

	flag "github.com/spf13/pflag"

	"gobackend/pkg/version"
)

// VersionValue ...
type VersionValue int

// Define some const.
const (
	VersionFalse VersionValue = 0
	VersionTrue  VersionValue = 1
	VersionRaw   VersionValue = 2
)

// version flag
const (
	versionFlagName      = "version"
	versionFlagShortName = "V"
)

const strRawVersion string = "raw"

// IsBoolFlag ...
func (v *VersionValue) IsBoolFlag() bool {
	return true
}

// Get ...
func (v *VersionValue) Get() interface{} {
	return v
}

// Set ...
func (v *VersionValue) Set(s string) error {
	if s == strRawVersion {
		*v = VersionRaw

		return nil
	}

	boolVal, err := strconv.ParseBool(s)
	if boolVal {
		*v = VersionTrue
	} else {
		*v = VersionFalse
	}

	return err
}

func (v *VersionValue) String() string {
	if *v == VersionRaw {
		return strRawVersion
	}

	return fmt.Sprintf("%v", bool(*v == VersionTrue))
}

// Type of the flag as required by the pflag.Value interface.
func (v *VersionValue) Type() string {
	return "version"
}

// VersionVar defines a flag with the specified name and usage string.
func VersionVar(p *VersionValue, name, shortName string, value VersionValue, usage string) {
	*p = value
	flag.VarP(p, name, shortName, usage)
	// "--version" will be treated as "--version=true"
	flag.Lookup(name).NoOptDefVal = "true"
}

// Version wraps the VersionVar function.
func Version(name, shortName string, value VersionValue, usage string) *VersionValue {
	p := new(VersionValue)
	VersionVar(p, name, shortName, value, usage)

	return p
}

var versionFlag = Version(versionFlagName, versionFlagShortName, VersionFalse, "Print version information and quit.")

// AddFlags registers this package's flags on arbitrary FlagSets, such that they point to the
// same value as the global flags.
func AddFlags(fs *flag.FlagSet) {
	fs.AddFlag(flag.Lookup(versionFlagName))
}

// PrintAndExitIfRequested will check if the -version flag was passed
// and, if so, print the version and exit.
func PrintAndExitIfRequested() {
	if *versionFlag == VersionRaw {
		fmt.Printf("%#v\n", version.Get())
		os.Exit(0)
	} else if *versionFlag == VersionTrue {
		fmt.Printf("%s\n", version.Get())
		os.Exit(0)
	}
}
