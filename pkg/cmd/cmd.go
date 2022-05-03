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
		if flag.Value.Type() == "string" {
			_ = flag.Value.Set(os.Getenv(flagName2EnvName(flag.Name))) // nolint: errcheck
		}
	})
}

func flagName2EnvName(name string) string {
	return strings.ToUpper(strings.Replace(name, "-", "_", -1))
}
