package config

import (
	"github.com/koding/multiconfig"
)

//Config - struct for project info.
type Config struct {
	//Port for connecting with GRPC server.
	PortGRPCServer string `default:":8001"`
	//Port for connecting with REST server.
	PortRESTServer string `default:":8000"`

	//Info for connecting with mysql database.
	Login    string `default:"root"`
	Passwd   string `default:"root"`
	Hostname string `default:"172.17.0.4"`
	Port     int    `default:"3306"`
	DBName   string `default:"busstation"`
}

//GetData - get data from config file(new config object).
func GetData() *Config {
	m := multiconfig.NewWithPath("config.toml")
	configStruct := new(Config)
	m.MustLoad(configStruct)
	return configStruct
}
