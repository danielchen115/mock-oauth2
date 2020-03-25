package oauth

import (
	"fmt"
	"github.com/spf13/viper"
)

type FieldSpec struct {
	Name     string      `mapstructure:"name"`
	Required bool        `mapstructure:"required"`
}

type Config struct {
	Import ImportConfig
}

type ImportConfig struct {
	Fields []FieldSpec `mapstructure:"fields"`
}

func LoadConfig(file string, dir string) (*Config, error) {
	viper.SetConfigFile(file)
	viper.AddConfigPath(dir)
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("config not read: %s", err)
		return nil, err
	}
	var c Config
	err = viper.Unmarshal(&c)
	if err != nil {
		fmt.Printf("config not loaded: %s", err)
		return nil, err
	}
	return &c, err
}