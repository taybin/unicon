// Package unicon provides tools for managing hierarcial configuration from multiple sources
package unicon

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cast"
	"strings"
	"time"
)

// Configurable is the main interface.  Also the hierarcial configuration (Config) implements it.
type Configurable interface {
	Get(string) interface{}
	GetString(key string) string
	GetBool(key string) bool
	GetInt(key string) int
	GetInt64(key string) int64
	GetFloat64(key string) float64
	GetTime(key string) time.Time
	GetDuration(key string) time.Duration

	// Set a variable, nil to reset key
	Set(string, interface{})
	// Reset the config data to passed data, if nothing is given set it to zero value
	Reset(...map[string]interface{})
	// Return a map of all variables
	All() map[string]interface{}
}

// ReadableConfig is a Configurable that can be loaded
type ReadableConfig interface {
	Configurable
	// Load the configuration
	Load() error
}

// WritableConfig is a Configurable that can be Loaded & Saved
type WritableConfig interface {
	ReadableConfig
	// Save configuration
	Save() error
}

// Config is a Configurable that can Use other Configurables thus build a hierarchy
type Config interface {
	WritableConfig
	// Use config as name, .Use("name") without the second parameter returns
	// the config previously added to the hierarchy with the name.
	// Use("name", Configurable) adds or replaces the configurable designated by "Name" in
	// the hierarchy
	Use(name string, config ...Configurable) Configurable
}

// Unicon is the Hierarchical Config that can be used to mount other configs that are searched for keys by Get
type Unicon struct {
	// Overrides, these are checked before Configs are iterated for key
	Configurable
	// named configurables, these are iterated if key is not found in Config
	Configs map[string]Configurable
	// Defaults configurable, if key is not found in the Configurable & Configurables in Config,
	//Defaults is checked for fallback values
	Defaults Configurable
	prefix   string
}

// Ensure Unicon implements Config
var _ Config = (*Unicon)(nil)

// NewConfig creates a new config that is by default backed by a MemoryConfig Configurable
// Takes optional initial configuration and an optional defaults
func NewConfig(initial Configurable, defaults ...Configurable) *Unicon {
	if initial == nil {
		initial = NewMemoryConfig()
	} else {
		LoadConfig(initial)
	}
	if len(defaults) == 0 {
		defaults = append(defaults, NewMemoryConfig())
	}

	return &Unicon{
		Configurable: initial,
		Configs:      make(map[string]Configurable),
		Defaults:     defaults[0],
		prefix:       "",
	}
}

// Unmarshal current configuration hierarchy into target using gonfig:
func (uni *Unicon) Unmarshal(target interface{}) error {
	err := mapstructure.WeakDecode(uni.All(), target)

	if err != nil {
		return err
	}

	return nil
}

// Reset resets all configs with the provided data, if no data is provided empties all stores
// Never touches the Defaults, to reset Defaults use Config.Defaults().Reset()
func (uni *Unicon) Reset(datas ...map[string]interface{}) {
	var data map[string]interface{}
	if len(datas) > 0 {
		data = datas[0]
	}
	for _, value := range uni.Configs {
		if data != nil {
			value.Reset(data)
		} else {
			value.Reset()
		}
	}
	uni.Configurable.Reset(data)
}

// Use config as named config and return an already set and loaded config
// mounts a new configuration in the hierarchy.
// conf.Use("global", NewUrlConfig("http://host.com/config..json")).
// conf.Use("local", NewFileConfig("./config.json"))
// err := conf.Load();.
// Then get variable from specific config.
// conf.Use("global").Get("key").
// or traverse the hierarchy and search for "key".
// conf.Get("key").
// conf.Use("name") returns a nil value for non existing config named "name".
func (uni *Unicon) Use(name string, config ...Configurable) Configurable {
	if uni.Configs == nil {
		uni.Configs = make(map[string]Configurable)
	}
	if len(config) == 0 {
		return uni.Configs[name]
	}
	uni.Configs[name] = config[0]
	LoadConfig(uni.Configs[name])
	return uni.Configs[name]
}

// Get gets the key from first store that it is found from, checks Defaults
func (uni *Unicon) Get(key string) interface{} {
	key = uni.prefixedKey(key)
	// override from out values
	if value := uni.Configurable.Get(key); value != nil {
		return value
	}
	// go through all in insert order until key is found
	for _, config := range uni.Configs {
		if value := config.Get(key); value != nil {
			return value
		}
	}
	// if not found check the defaults as fallback
	if value := uni.Defaults.Get(key); value != nil {
		return value
	}

	return nil
}

func (uni *Unicon) GetString(key string) string {
	return cast.ToString(uni.Get(key))
}

func (uni *Unicon) GetBool(key string) bool {
	return cast.ToBool(uni.Get(key))
}

func (uni *Unicon) GetInt(key string) int {
	return cast.ToInt(uni.Get(key))
}

func (uni *Unicon) GetInt64(key string) int64 {
	return cast.ToInt64(uni.Get(key))
}

func (uni *Unicon) GetFloat64(key string) float64 {
	return cast.ToFloat64(uni.Get(key))
}

func (uni *Unicon) GetTime(key string) time.Time {
	return cast.ToTime(uni.Get(key))
}

func (uni *Unicon) GetDuration(key string) time.Duration {
	return cast.ToDuration(uni.Get(key))
}

func (uni *Unicon) Set(key string, value interface{}) {
	key = uni.prefixedKey(key)
	uni.Configurable.Set(key, value)
}

// SaveConfig saves if is of type WritableConfig, otherwise does nothing.
func SaveConfig(config Configurable) error {
	switch t := config.(type) {
	case WritableConfig:

		if err := t.Save(); err != nil {
			return err
		}
	}
	return nil
}

// Save saves all mounted configurations in the hierarchy that implement the WritableConfig interface
func (uni *Unicon) Save() error {
	for _, config := range uni.Configs {
		if err := SaveConfig(config); err != nil {
			return err
		}
	}
	return SaveConfig(uni.Configurable)
}

// LoadConfig loads a config if it is of type ReadableConfig, otherwise does nothing.
func LoadConfig(config Configurable) error {
	switch t := config.(type) {
	case ReadableConfig:
		if err := t.Load(); err != nil {
			return err
		}
	}
	return nil
}

// Load calls Configurable.Load() on all Configurable objects in the hierarchy.
func (uni *Unicon) Load() error {
	LoadConfig(uni.Configurable)
	LoadConfig(uni.Defaults)
	for _, config := range uni.Configs {
		LoadConfig(config)
	}
	return nil
}

// All returns a map of data from all Configurables in use
// the first found instance of variable found is provided.
// Config.Use("a", NewMemoryConfig()).
// Config.Use("b", NewMemoryConfig()).
// Config.Use("a").Set("a","1").
// Config.Set("b").Set("a","2").
// then.
// Config.All()["a"] == "1".
// Config.Get("a") == "1".
// Config.Use("b".).Get("a") == "2".
func (uni *Unicon) All() map[string]interface{} {
	values := make(map[string]interface{})
	// put defaults in values
	for key, value := range uni.Defaults.All() {
		if values[key] == nil {
			values[key] = value
		}
	}
	// put config values on top of them
	for _, config := range uni.Configs {
		for key, value := range config.All() {
			if values[key] == nil {
				values[key] = value
			}
		}
	}
	// put overrides from uni on top of all
	for key, value := range uni.Configurable.All() {
		if values[key] == nil {
			values[key] = value
		}
	}
	return values
}

// trimsplit slices s into all substrings separated by sep and returns a
// slice of the substrings between the separator with all leading and trailing
// white space removed, as defined by Unicode.
func trimsplit(s, sep string) []string {
	trimmed := strings.Split(s, sep)
	for i := range trimmed {
		trimmed[i] = strings.TrimSpace(trimmed[i])
	}
	return trimmed
}

// Sub returns a new Unicon but with the namespace prepended to Gets/Sets/Subs
// behind the scenes
func (uni *Unicon) Sub(ns string) *Unicon {
	sub := NewConfig(uni, uni.Defaults)
	sub.prefix = uni.prefixedKey(ns)
	return sub
}

func (uni *Unicon) prefixedKey(key string) string {
	if uni.prefix != "" {
		return strings.Join([]string{uni.prefix, key}, ".")
	}
	return key
}

// Debug prints out simple list of keys as returned by All()
func (uni *Unicon) Debug() {
	for key := range uni.All() {
		fmt.Printf("%s = %s\n", key, uni.GetString(key))
	}
}
