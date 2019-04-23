package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/jessevdk/go-flags"
	"gopkg.in/yaml.v2"
)

//Config - struct for project info.
type Config struct {
	//Port for connecting with GRPC server.
	PortGRPCServer string `short:"g" default:":8001"`
	//Port for connecting with REST server.
	PortRESTServer string `short:"r" default:":8000"`
	ConfigFile     string `short:"f" long:"configfile" description:"File with config"`
	Port           string `short:"d" long:"dbport" default:":3306"`
	Login          string `long:"login" default:"root"`
	Passwd         string `long:"password" default:"root"`
	Hostname       string `long:"hostname" default:"172.18.0.2"`
	DBName         string `long:"dbname" default:"busstation"`
}

//Parse works with command arguments.
func (c *Config) Parse() error {
	parser := flags.NewParser(c, flags.Default|flags.IgnoreUnknown)
	_, err := parser.Parse()
	if err != nil {
		return fmt.Errorf("parse error: %v", err)
	}
	if c.ConfigFile != "" {
		err = c.LoadOptionsFromFile()
		if err != nil {
			log.Printf("cannot read settings from file: %v", err)
		}
	}
	return nil
}

//GetData - get data from config file(new config object).
func (c *Config) LoadOptionsFromFile() error {
	data, err := ioutil.ReadFile(c.ConfigFile)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, c)
}
