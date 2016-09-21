package unicon

import (
	"flag"
	"github.com/spf13/pflag"
	"os"
	"strings"
)

// PflagConfig can be used to read arguments in the posix flag style
// into the underlaying Configurable
type PflagConfig struct {
	Configurable
	Prefix     string
	namespaces []string
}

// NewPflagConfig creates a new PflagConfig and returns it as a ReadableConfig
func NewPflagConfig(prefix string, namespaces ...string) ReadableConfig {
	cfg := &PflagConfig{
		Configurable: NewMemoryConfig(),
		Prefix:       prefix,
		namespaces:   nsSlice(namespaces),
	}
	return cfg
}

// Load loads all the variables from argv to the underlaying Configurable.
// If a Prefix is provided for PflagConfig then keys are imported with the
// Prefix removed so --test.asd=1 with Prefix 'test.' imports "asd" with
// value of 1
func (pc *PflagConfig) Load() (err error) {
	flagset := pflag.NewFlagSet("arguments", pflag.ContinueOnError)
	flagset.Parse(os.Args)

	flagset.VisitAll(func(f *pflag.Flag) {
		name := f.Name
		if pc.Prefix != "" && strings.HasPrefix(f.Name, pc.Prefix) {
			name = strings.Replace(name, pc.Prefix, "", 1)
		}

		var value interface{}
		if getter, ok := f.Value.(flag.Getter); ok {
			value = getter.Get().(string)
		} else {
			value = f.Value.String()
		}

		name = namespaceKey(name, pc.namespaces)
		pc.Set(name, value)
	})
	return nil
}
