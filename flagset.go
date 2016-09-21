package unicon

import (
	"flag"
	"github.com/spf13/pflag"
	"strings"
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
	// put in lowercase
	var lowered []string
	for _, ns := range namespaces {
		lowered = append(lowered, strings.ToLower(ns))
	}

	cfg := &FlagSetConfig{
		Configurable: NewMemoryConfig(),
		Prefix:       prefix,
		namespaces:   lowered,
		fs:           fs,
	}
	return cfg
}

// Load loads all the variables from argv to the underlaying Configurable.
// If a Prefix is provided for FlagSetConfig then keys are imported with the Prefix removed
// so --test.asd=1 with Prefix 'test.' imports "asd" with value of 1
func (pc *FlagSetConfig) Load() (err error) {
	pc.fs.VisitAll(func(f *pflag.Flag) {
		name := f.Name
		if pc.Prefix != "" && strings.HasPrefix(f.Name, pc.Prefix) {
			name = strings.Replace(name, pc.Prefix, "", 1)
		}
		var value interface{}
		if getter, ok := f.Value.(flag.Getter); ok {
			value = getter.Get().(string)
			pc.Set(name, getter.Get().(string))
		} else {
			value = f.Value.String()
		}

		name = namespaceKey(name, pc.namespaces)
		pc.Set(name, value)
	})
	return nil
}
