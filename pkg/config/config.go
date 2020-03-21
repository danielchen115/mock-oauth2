package config

import (
	"encoding/json"
	"os"
)

type Field struct {
	Name     string      `yaml:"name"`
	Required bool        `yaml:"required"`
}

type Config struct {
	Fields []Field `yaml:"fields"`
}

func Load(file string) (*Config, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(f)
	var c Config
	err = decoder.Decode(&c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}