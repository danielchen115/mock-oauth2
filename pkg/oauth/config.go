package oauth

import (
	"gopkg.in/yaml.v2"
	"os"
)

type FieldSpec struct {
	Name     string      `yaml:"name"`
	Required bool        `yaml:"required"`
}

type Config struct {
	Fields []FieldSpec `yaml:"fields"`
}

func LoadConfig(file string) (*Config, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	decoder := yaml.NewDecoder(f)
	var c Config
	err = decoder.Decode(&c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}