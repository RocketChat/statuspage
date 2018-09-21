package config

import (
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"
)

var Config *config

type config struct {
	DataPath string          `yaml:"dataPath"`
	Services []serviceConfig `yaml:"services"`
}

type serviceConfig struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

func (c *config) Load(filePath string) error {

	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
		return err
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
		return err
	}

	return nil
}

func Load(filePath string) error {
	Config = new(config)
	err := Config.Load(filePath)
	if err != nil {
		return err
	}

	return nil
}
