package unicon

import (
	"os"
	"strings"
)

// EnvConfig can be used to read values from the environment
// into the underlaying Configurable
type EnvConfig struct {
	Configurable
	Prefix     string
	namespaces []string
}

// NewEnvConfig creates a new Env config backed by a memory config
func NewEnvConfig(prefix string, namespaces ...string) ReadableConfig {
	// put in lowercase
	var lowered []string
	for _, ns := range namespaces {
		lowered = append(lowered, strings.ToLower(ns))
	}

	cfg := &EnvConfig{
		Configurable: NewMemoryConfig(),
		Prefix:       prefix,
		namespaces:   lowered,
	}
	cfg.Load()
	return cfg
}

// Load loads the data from os.Environ() to the underlaying Configurable.
// if a Prefix is set then variables are imported with self.Prefix removed from the name
// so MYAPP_test=1 exported in env and read from ENV by EnvConfig{Prefix:"MYAPP_"} can be found from
// EnvConfig.Get("test")
// If namespaces are declared, POSTGRESQL_HOST becomes postgresql.host
func (ec *EnvConfig) Load() (err error) {
	env := os.Environ()
	for _, pair := range env {
		kv := strings.Split(pair, "=")
		if kv != nil && len(kv) >= 2 {
			name := strings.Replace(kv[0], ec.Prefix, "", 1)
			name = namespaceKey(name, ec.namespaces)
			ec.Set(name, kv[1])
		}
	}
	return nil
}
