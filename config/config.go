package config

import (
	"errors"
	"io/ioutil"
	"log"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

//Config contains a reference to the configuration
var Config *config

type config struct {
	DataPath string          `yaml:"dataPath" json:"dataPath"`
	Services []serviceConfig `yaml:"services" json:"services"`
}

type serviceConfig struct {
	Name        string `yaml:"name" json:"name"`
	Description string `yaml:"description" json:"description"`
}

func (c *config) Load(filePath string) error {
	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
		return err
	}

	if err = yaml.Unmarshal(yamlFile, c); err != nil {
		log.Fatalf("Unmarshal: %v", err)
		return err
	}

	return nil
}

func (c *config) VerifySettings() error {
	if c.DataPath == "" {
		return errors.New("invalid dataPath, can not be empty")
	}

	if !strings.HasSuffix(c.DataPath, "/") {
		return errors.New("dataPath must end with '/'")
	}

	return nil
}

//Load tries to load the configuration file
func Load(filePath string) error {
	Config = new(config)

	if err := Config.Load(filePath); err != nil {
		return err
	}

	return Config.VerifySettings()
}
