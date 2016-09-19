// package unicon provides tools for managing hierarcial configuration from multiple sources
package unicon

import (
	"fmt"
	"github.com/spf13/cast"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// The main Configurable interface Also the hierarcial configuration (Config) implements it.
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

// A Configurable that can be loaded
type ReadableConfig interface {
	Configurable
	// Load the configuration
	Load() error
}

// A Configurable that can be Loaded & Saved
type WritableConfig interface {
	ReadableConfig
	// Save configuration
	Save() error
}

// a Configurable that can Use other Configurables thus build a hierarchy
type Config interface {
	WritableConfig
	// Use config as name, .Use("name") without the second parameter returns
	// the config previously added to the hierarchy with the name.
	// Use("name", Configurable) adds or replaces the configurable designated by "Name" in
	// the hierarchy
	Use(name string, config ...Configurable) Configurable
}

// The Hierarchical Config that can be used to mount other configs that are searched for keys by Get
type Unicon struct {
	// Overrides, these are checked before Configs are iterated for key
	Configurable
	// named configurables, these are iterated if key is not found in Config
	Configs map[string]Configurable
	// Defaults configurable, if key is not found in the Configurable & Configurables in Config,
	//Defaults is checked for fallback values
	Defaults Configurable
}

// Ensure Unicon implements Config
var _ Config = (*Unicon)(nil)

// Creates a new config that is by default backed by a MemoryConfig Configurable
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
	}
}

// Marshal current configuration hierarchy into target using gonfig:
func (self *Unicon) Marshal(target interface{}) error {
	value := reflect.Indirect(reflect.ValueOf(target))
	typ := value.Type()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		v := strings.TrimSpace(self.GetString(field.Tag.Get("unicon")))
		if v == "" {
			continue
		}

		// Set the appropriate type.
		switch field.Type.Kind() {
		case reflect.Bool:
			value.Field(i).SetBool(v != "0" && v != "false")
		case reflect.Int:
			newValue, err := strconv.ParseInt(v, 10, 0)
			if err != nil {
				return fmt.Errorf("Parse error: %s: %s", field.Tag.Get("env"), err)
			}
			value.Field(i).SetInt(newValue)
		case reflect.String:
			value.Field(i).SetString(v)
		case reflect.Slice:
			value.Field(i).Set(reflect.ValueOf(trimsplit(v, ",")))
		}
	}
	return nil
}

// Resets all configs with the provided data, if no data is provided empties all stores
// Never touches the Defaults, to reset Defaults use Config.Defaults().Reset()
func (self *Unicon) Reset(datas ...map[string]interface{}) {
	var data map[string]interface{}
	if len(datas) > 0 {
		data = datas[0]
	}
	for _, value := range self.Configs {
		if data != nil {
			value.Reset(data)
		} else {
			value.Reset()
		}
	}
	self.Configurable.Reset(data)
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
func (self *Unicon) Use(name string, config ...Configurable) Configurable {
	if self.Configs == nil {
		self.Configs = make(map[string]Configurable)
	}
	if len(config) == 0 {
		return self.Configs[name]
	}
	self.Configs[name] = config[0]
	LoadConfig(self.Configs[name])
	return self.Configs[name]
}

func (self *Unicon) GetString(key string) string {
	return cast.ToString(self.Get(key))
}

func (self *Unicon) GetBool(key string) bool {
	return cast.ToBool(self.Get(key))
}

func (self *Unicon) GetInt(key string) int {
	return cast.ToInt(self.Get(key))
}

func (self *Unicon) GetInt64(key string) int64 {
	return cast.ToInt64(self.Get(key))
}

func (self *Unicon) GetFloat64(key string) float64 {
	return cast.ToFloat64(self.Get(key))
}

func (self *Unicon) GetTime(key string) time.Time {
	return cast.ToTime(self.Get(key))
}

func (self *Unicon) GetDuration(key string) time.Duration {
	return cast.ToDuration(self.Get(key))
}

// Gets the key from first store that it is found from, checks Defaults
func (self *Unicon) Get(key string) interface{} {
	// override from out values
	if value := self.Configurable.Get(key); value != nil {
		return value
	}
	// go through all in insert order until key is found
	for _, config := range self.Configs {
		if value := config.Get(key); value != nil {
			return value
		}
	}
	// if not found check the defaults as fallback
	if value := self.Defaults.Get(key); value != nil {
		return value
	}

	return nil
}

func (self *Unicon) Set(key string, value interface{}) {
	self.Configurable.Set(key, value)
}

// Save config it is of type WritableConfig, otherwise does nothing.
func SaveConfig(config Configurable) error {
	switch t := config.(type) {
	case WritableConfig:

		if err := t.Save(); err != nil {
			return err
		}
	}
	return nil
}

// Saves all mounted configurations in the hierarchy that implement the WritableConfig interface
func (self *Unicon) Save() error {
	for _, config := range self.Configs {
		if err := SaveConfig(config); err != nil {
			return err
		}
	}
	return SaveConfig(self.Configurable)
}

// Load config it is of type ReadableConfig, otherwise does nothing.
func LoadConfig(config Configurable) error {
	switch t := config.(type) {
	case ReadableConfig:
		if err := t.Load(); err != nil {
			return err
		}
	}
	return nil
}

// calls Configurable.Load() on all Configurable objects in the hierarchy.
func (self *Unicon) Load() error {
	LoadConfig(self.Configurable)
	LoadConfig(self.Defaults)
	for _, config := range self.Configs {
		LoadConfig(config)
	}
	return nil
}

// Returns a map of data from all Configurables in use
// the first found instance of variable found is provided.
// Config.Use("a", NewMemoryConfig()).
// Config.Use("b", NewMemoryConfig()).
// Config.Use("a").Set("a","1").
// Config.Set("b").Set("a","2").
// then.
// Config.All()["a"] == "1".
// Config.Get("a") == "1".
// Config.Use("b".).Get("a") == "2".
func (self *Unicon) All() map[string]interface{} {
	values := make(map[string]interface{})
	// put defaults in values
	for key, value := range self.Defaults.All() {
		if values[key] == nil {
			values[key] = value
		}
	}
	// put config values on top of them
	for _, config := range self.Configs {
		for key, value := range config.All() {
			if values[key] == nil {
				values[key] = value
			}
		}
	}
	// put overrides from self on top of all
	for key, value := range self.Configurable.All() {
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
