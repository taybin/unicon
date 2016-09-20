package unicon

// MemoryConfig is a simple abstraction to map[]interface{} for in process memory backed configuration
// only implements Configurable use JsonConfig to save/load if needed
type MemoryConfig struct {
	data map[string]interface{}
}

// NewMemoryConfig returns a new memory backed Configurable
// The most basic Configurable simply backed by a map[string]interface{}
func NewMemoryConfig() *MemoryConfig {
	cfg := &MemoryConfig{
		data: make(map[string]interface{}),
	}
	return cfg
}

func (mem *MemoryConfig) init() {
	mem.data = make(map[string]interface{})
}

// Reset if no arguments are provided Reset() re-creates the underlaying map
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
	return mem.data[key]
}

// All returns all keys
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
	mem.data[key] = value
}
