package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	DatabaseHost	 string `yaml:"db_host"`
	DatabasePort	 int    `yaml:"db_port"`
	DatabaseUser	 string `yaml:"db_user"`
	DatabasePassword string `yaml:"db_password"`
	DatabaseName	 string `yaml:"db_name"`
	ExecuteAutomatic bool   `yaml:"execute_automatic"`
}

func GetConfigValues(cfgPath string) (*Config, error) {
	cfgFile := cfgPath
	cfgByte, err := os.ReadFile(cfgFile)
	if err != nil {
		return nil, err
	}
	var config Config
	err = yaml.Unmarshal(cfgByte, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
