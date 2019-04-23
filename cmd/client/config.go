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
	//Port for client work.
	PortClient string `short:"c" default:":8080"`
	ConfigFile string `short:"f" long:"configfile" description:"File with config"`
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

//GetData - get data from config file.
func (c *Config) LoadOptionsFromFile() error {
	data, err := ioutil.ReadFile(c.ConfigFile)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, c)
}
