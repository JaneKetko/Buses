package config

import (
	"github.com/koding/multiconfig"
)

//Config - struct for project info.
type Config struct {
	PortServer int    `default:"8000"`
	Login      string `default:"root"`
	Passwd     string `default:"root"`
	Hostname   string `default:"172.17.0.2"`
	Port       int    `default:"3306"`
	DBName     string `default:"busstation"`
}

//GetData - get data from config file(new config object).
func GetData() *Config {
	m := multiconfig.NewWithPath("config.toml")
	configStruct := new(Config)
	m.MustLoad(configStruct)
	return configStruct
}
