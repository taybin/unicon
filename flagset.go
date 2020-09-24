package unicon

import (
	"flag"
	"strings"

	"github.com/spf13/pflag"
)

// FlagSetConfig can be used to read arguments in the posix flag style
// into the underlaying Configurable
type FlagSetConfig struct {
	Configurable
	Prefix     string
	namespaces []string
	fs         *pflag.FlagSet
}

// NewFlagSetConfig creates a new FlagSetConfig and returns it as a
// ReadableConfig
func NewFlagSetConfig(fs *pflag.FlagSet, prefix string, namespaces ...string) ReadableConfig {
	cfg := &FlagSetConfig{
		Configurable: NewMemoryConfig(),
		Prefix:       prefix,
		namespaces:   nsSlice(namespaces),
		fs:           fs,
	}
	return cfg
}

// Load loads all the variables from argv to the underlaying Configurable.
// If a Prefix is provided for FlagSetConfig then keys are imported with the
// Prefix removed so --test.asd=1 with Prefix 'test.' imports "asd" with
// value of 1
func (fsc *FlagSetConfig) Load() (err error) {
	fsc.fs.VisitAll(func(f *pflag.Flag) {
		name := f.Name
		if fsc.Prefix != "" && strings.HasPrefix(f.Name, fsc.Prefix) {
			name = strings.Replace(name, fsc.Prefix, "", 1)
		}
		var value interface{}
		if getter, ok := f.Value.(flag.Getter); ok {
			value = getter.Get().(string)
		} else {
			value = f.Value.String()
		}

		name = namespaceKey(name, fsc.namespaces)
		fsc.Set(name, value)
	})
	return nil
}
