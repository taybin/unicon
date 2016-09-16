package gonfig

import (
	"flag"
	"github.com/ogier/pflag"
	"log"
	"os"
	"strings"
)

// the Pflag configurable.
type PflagConfig struct {
	Configurable
	Prefix string
}

// Creates a new PflagConfig and returns it as a ReadableConfig
func NewPflagConfig(prefix string) ReadableConfig {
	cfg := &PflagConfig{NewMemoryConfig(), prefix}
	return cfg
}

// Loads all the variables from argv to the underlaying Configurable.
// If a Prefix is provided for PflagConfig then keys are imported with the Prefix removed
// so --test.asd=1 with Prefix 'test.' imports "asd" with value of 1
func (self *PflagConfig) Load() (err error) {
	flagset := pflag.NewFlagSet("arguments", pflag.ContinueOnError)
	flagset.Parse(os.Args)

	flagset.VisitAll(func(f *pflag.Flag) {
		name := f.Name
		log.Println(f.Name)
		if self.Prefix != "" && strings.HasPrefix(f.Name, self.Prefix) {
			name = strings.Replace(name, self.Prefix, "", 1)
		}
		if getter, ok := f.Value.(flag.Getter); ok {
			self.Set(name, getter.Get().(string))
		}
	})
	return nil
}
