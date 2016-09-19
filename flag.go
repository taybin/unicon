package unicon

import (
	"flag"
	"os"
	"strings"
)

// FlagConfig can be used to read arguments in the golang flag style
// into the underlaying Configurable
type FlagConfig struct {
	Configurable
	Prefix     string
	namespaces []string
}

// NewFlagConfig creates a new FlagConfig and returns it as a ReadableConfig
func NewFlagConfig(prefix string, namespaces ...string) ReadableConfig {
	// put in lowercase
	var lowered []string
	for _, ns := range namespaces {
		lowered = append(lowered, strings.ToLower(ns))
	}

	cfg := &FlagConfig{
		Configurable: NewMemoryConfig(),
		Prefix:       prefix,
		namespaces:   lowered,
	}
	return cfg
}

// Load loads all the variables from argv to the underlaying Configurable.
// If a Prefix is provided for FlagConfig then keys are imported with the Prefix removed
// so --test.asd=1 with Prefix 'test.' imports "asd" with value of 1
func (pc *FlagConfig) Load() (err error) {
	flagset := flag.NewFlagSet("arguments", flag.ContinueOnError)
	flagset.Parse(os.Args)

	flagset.VisitAll(func(f *flag.Flag) {
		name := f.Name
		if pc.Prefix != "" && strings.HasPrefix(f.Name, pc.Prefix) {
			name = strings.Replace(name, pc.Prefix, "", 1)
		}
		if getter, ok := f.Value.(flag.Getter); ok {
			pc.Set(name, getter.Get().(string))
		} else {
			pc.Set(name, f.Value.String())
		}
	})
	return nil
}
