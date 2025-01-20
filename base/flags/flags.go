package flags

import (
	"os"

	"github.com/spf13/pflag"
)

func NewFlagSet() *pflag.FlagSet {
	// Create a new flagset for global configs.
	flagset := pflag.NewFlagSet("global config", pflag.ExitOnError)

	// Add flags to the flagset.
	flagset.BoolP("development", "d", false, "toggle development mode")
	flagset.Uint16P("port", "p", 8080, "the listening port of the http web server")

	// Parse the flagset.
	flagset.Parse(os.Args[1:])

	return flagset
}
