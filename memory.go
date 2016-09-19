package unicon

import (
	"github.com/spf13/cast"
	"strings"
	"time"
)

// MemoryConfig is a simple abstraction to map[]interface{} for in process memory backed configuration
// only implements Configurable use JsonConfig to save/load if needed
type MemoryConfig struct {
	data   map[string]interface{}
	prefix string
}

// NewMemoryConfig returns a new memory backed Configurable
// The most basic Configurable simply backed by a map[string]interface{}
func NewMemoryConfig() *MemoryConfig {
	cfg := &MemoryConfig{
		data:   make(map[string]interface{}),
		prefix: "",
	}
	return cfg
}

func (mem *MemoryConfig) init() {
	mem.data = make(map[string]interface{})
}

// if no arguments are proced Reset() re-creates the underlaying map
func (mem *MemoryConfig) Reset(datas ...map[string]interface{}) {
	if len(datas) >= 1 {
		mem.data = datas[0]
	} else {
		mem.data = make(map[string]interface{})
	}
}

// Get key from map
func (mem *MemoryConfig) Get(key string) interface{} {
	if mem.data == nil {
		mem.init()
	}
	key = mem.prefixedKey(key)
	return mem.data[key]
}

func (mem *MemoryConfig) GetString(key string) string {
	return cast.ToString(mem.Get(key))
}

func (mem *MemoryConfig) GetBool(key string) bool {
	return cast.ToBool(mem.Get(key))
}

func (mem *MemoryConfig) GetInt(key string) int {
	return cast.ToInt(mem.Get(key))
}

func (mem *MemoryConfig) GetInt64(key string) int64 {
	return cast.ToInt64(mem.Get(key))
}

func (mem *MemoryConfig) GetFloat64(key string) float64 {
	return cast.ToFloat64(mem.Get(key))
}

func (mem *MemoryConfig) GetTime(key string) time.Time {
	return cast.ToTime(mem.Get(key))
}

func (mem *MemoryConfig) GetDuration(key string) time.Duration {
	return cast.ToDuration(mem.Get(key))
}

// get all keys
func (mem *MemoryConfig) All() map[string]interface{} {
	if mem.data == nil {
		mem.init()
	}
	return mem.data
}

// Set a key to value
func (mem *MemoryConfig) Set(key string, value interface{}) {
	if mem.data == nil {
		mem.init()
	}
	key = mem.prefixedKey(key)
	mem.data[key] = value
}

func (mem *MemoryConfig) setPrefix(ns string) {
	mem.prefix = ns
}

func (mem *MemoryConfig) prefixedKey(key string) string {
	if mem.prefix != "" {
		return strings.Join([]string{mem.prefix, key}, ".")
	}
	return key
}
