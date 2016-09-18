package gonfig

import (
	"github.com/spf13/cast"
	"time"
)

// MemoryConfig is a simple abstraction to map[]interface{} for in process memory backed configuration
// only implements Configurable use JsonConfig to save/load if needed
type MemoryConfig struct {
	data map[string]interface{}
}

// Returns a new memory backed Configurable
// The most basic Configurable simply backed by a map[string]interface{}
func NewMemoryConfig() *MemoryConfig {
	cfg := &MemoryConfig{make(map[string]interface{})}
	cfg.init()
	return cfg
}

func (self *MemoryConfig) init() {
	self.data = make(map[string]interface{})
}

// if no arguments are proced Reset() re-creates the underlaying map
func (self *MemoryConfig) Reset(datas ...map[string]interface{}) {
	if len(datas) >= 1 {
		self.data = datas[0]
	} else {
		self.data = make(map[string]interface{})
	}
}

// Get key from map
func (self *MemoryConfig) Get(key string) interface{} {
	if self.data == nil {
		self.init()
	}
	return self.data[key]
}

func (self *MemoryConfig) GetString(key string) string {
	return cast.ToString(self.Get(key))
}

func (self *MemoryConfig) GetBool(key string) bool {
	return cast.ToBool(self.Get(key))
}

func (self *MemoryConfig) GetInt(key string) int {
	return cast.ToInt(self.Get(key))
}

func (self *MemoryConfig) GetInt64(key string) int64 {
	return cast.ToInt64(self.Get(key))
}

func (self *MemoryConfig) GetFloat64(key string) float64 {
	return cast.ToFloat64(self.Get(key))
}

func (self *MemoryConfig) GetTime(key string) time.Time {
	return cast.ToTime(self.Get(key))
}

func (self *MemoryConfig) GetDuration(key string) time.Duration {
	return cast.ToDuration(self.Get(key))
}

// get all keys
func (self *MemoryConfig) All() map[string]interface{} {
	if self.data == nil {
		self.init()
	}
	return self.data
}

// Set a key to value
func (self *MemoryConfig) Set(key string, value interface{}) {
	if self.data == nil {
		self.init()
	}
	self.data[key] = value
}
