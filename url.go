package unicon

import (
	"io/ioutil"
	"net/http"
)

type URLConfig struct {
	Configurable
	url string
}

// Returns a new Configurable backed by JSON at url
func NewURLConfig(url string) ReadableConfig {
	return &URLConfig{NewMemoryConfig(), url}
}

func (uc *URLConfig) Load() error {
	resp, err := http.Get(uc.url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	out, err := unmarshalJSON(body)
	if err != nil {
		return err
	}
	uc.Reset(out)
	return nil
}
