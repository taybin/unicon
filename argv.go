package unicon

import (
	"flag"
	"os"
	"strings"
)

// ArgvConfig is the Argv configurable.
type ArgvConfig struct {
	Configurable
	Prefix     string
}

// Creates a new ArgvConfig and returns it as a ReadableConfig
func NewArgvConfig(prefix string) ReadableConfig {
	cfg := &ArgvConfig{NewMemoryConfig(), prefix}
	return cfg
}

// Load loads all the variables from argv to the underlaying Configurable.
// If a Prefix is provided for ArgvConfig then keys are imported with the
// Prefix removed so --test.asd=1 with Prefix 'test.' imports "asd" with
// value of 1
func (ac *ArgvConfig) Load() (err error) {
	flagset := flag.NewFlagSet("arguments", flag.ContinueOnError)
	flagset.Parse(os.Args)

	flagset.VisitAll(func(f *flag.Flag) {
		name := f.Name
		if ac.Prefix != "" && strings.HasPrefix(f.Name, ac.Prefix) {
			name = strings.Replace(name, ac.Prefix, "", 1)
		}
		if getter, ok := f.Value.(flag.Getter); ok {
			ac.Set(name, getter.Get().(string))
		}
	})
	return nil
}
