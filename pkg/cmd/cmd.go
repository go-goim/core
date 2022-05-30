package cmd

import (
	"os"
	"strings"

	"github.com/spf13/pflag"
)

// GlobalFlagSet flags
var GlobalFlagSet = pflag.NewFlagSet("global", pflag.ExitOnError)

// ParseFlags parses the command line flags.
func ParseFlags() error {
	SetFlagsFromEnv()
	return GlobalFlagSet.Parse(os.Args[1:])
}

func SetFlagsFromEnv() {
	GlobalFlagSet.VisitAll(func(flag *pflag.Flag) {
		val := os.Getenv(flagName2EnvName(flag.Name))
		if val != "" {
			flag.Value.Set(val) // nolint: errcheck
		}
	})
}

func flagName2EnvName(name string) string {
	return strings.ToUpper(strings.ReplaceAll(name, "-", "_"))
}
