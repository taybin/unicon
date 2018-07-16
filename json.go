package unicon

import (
	"encoding/json"
	"io/ioutil"
)

// JSONConfig is the json configurable
type JSONConfig struct {
	Configurable
	Path string
}

func unmarshalJSON(bytes []byte) (map[string]interface{}, error) {
	out := make(map[string]interface{})
	if err := json.Unmarshal(bytes, &out); err != nil {
		return nil, err
	}

	output := make(map[string]interface{})
	unmarshalMap(out, "", output)

	return output, nil
}

// NewJSONConfig returns a new WritableConfig backed by a json file at path.
// The file does not need to exist, if it does not exist the first Save call
// will create it.
func NewJSONConfig(path string, cfg ...Configurable) *JSONConfig {
	if len(cfg) == 0 {
		cfg = append(cfg, NewMemoryConfig())
	}
	LoadConfig(cfg[0])
	conf := &JSONConfig{cfg[0], path}
	LoadConfig(conf)
	return conf
}

// Load attempts to load the json configuration at JSONConfig.Path
// and Set them into the underlaying Configurable
func (jc *JSONConfig) Load() (err error) {
	var data []byte
	if data, err = ioutil.ReadFile(jc.Path); err != nil {
		return
	}
	out, err := unmarshalJSON(data)
	if err != nil {
		return
	}

	jc.Configurable.Reset(out)
	return
}

// Save attempts to save the configuration from the underlaying Configurable
// to json file at JSONConfig.Path
func (jc *JSONConfig) Save() (err error) {
	b, err := json.Marshal(jc.Configurable.All())
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(jc.Path, b, 0600); err != nil {
		return err
	}

	return nil
}
