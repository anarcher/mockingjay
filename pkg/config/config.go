package config

import (
	"gopkg.in/yaml.v2"

	"io/ioutil"
)

type Config struct {
	Proxies []Proxy `yaml:proxies`
}

type Proxy struct {
	Target string `yaml:target,omitempty` // If target is empty, just check filter for proxing or not.
	Filter Filter `yaml:filter`
}

type Filter struct {
	Form   map[string]string `yaml:form`
	Method string            `yaml:method,omitempty`
}

func ReadConfigBytes(bs []byte) (Config, error) {
	cfg := Config{}
	err := yaml.Unmarshal(bs, &cfg)
	return cfg, err
}

func ReadConfigFile(filename string) (Config, error) {
	bs, err := ioutil.ReadFile(filename)
	if err != nil {
		return Config{}, err
	}

	return ReadConfigBytes(bs)
}
