package config

import (
	"errors"
	"io/ioutil"
	"log"
	"strings"

	"github.com/caarlos0/env"
	"github.com/gin-gonic/gin"

	yaml "gopkg.in/yaml.v2"
)

//Config contains a reference to the configuration
var Config *config

type config struct {
	HTTP      httpConfig      `yaml:"http" json:"http"`
	DataPath  string          `yaml:"dataPath" json:"dataPath"`
	AuthToken string          `yaml:"authToken" json:"-" env:"AUTH_TOKEN"`
	Website   websiteConfig   `yaml:"website" json:"website"`
	Services  []serviceConfig `yaml:"services" json:"services"`
	Regions   []regionConfig  `yaml:"regions" json:"regions"`
	Twitter   twitterConfig   `yaml:"twitter" json:"twitter"`
}

type httpConfig struct {
	Port int `yaml:"port" json:"port"`
}

type websiteConfig struct {
	HeaderBgColor   string `yaml:"headerBgColor" json:"headerBgColor"`
	Title           string `yaml:"title" json:"title"`
	CacheBreaker    string `yaml:"cacheBreaker" json:"cacheBreaker"`
	DaysToAggregate int    `yaml:"daysToAggregate" json:"daysToAggregate"`
}

type serviceConfig struct {
	Name        string `yaml:"name" json:"name"`
	Description string `yaml:"description" json:"description"`
}

type regionConfig struct {
	Name        string `yaml:"name" json:"name"`
	Description string `yaml:"description" json:"description"`
	RegionCode  string `yaml:"regionCode" json:"regionCode"`
	ServiceName string `yaml:"serviceName" json:"serviceName"`
}

type twitterConfig struct {
	Enabled        bool   `yaml:"enabled" json:"enabled"`
	ConsumerKey    string `yaml:"consumerKey" json:"consumerKey"`
	ConsumerSecret string `yaml:"consumerSecret" json:"consumerSecret" env:"TWITTER_CONSUMER_SECRET"`
	AccessToken    string `yaml:"accessToken" json:"accessToken" env:"TWITTER_ACCESS_TOKEN"`
	AccessSecret   string `yaml:"accessSecret" json:"accessSecret" env:"TWITTER_ACCESS_SECRET"`
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

func (c *config) HttpHandler(gc *gin.Context) {
	gc.JSON(200, Config)
}

//Load tries to load the configuration file
func Load(filePath string) error {
	Config = new(config)

	if err := Config.Load(filePath); err != nil {
		return err
	}

	return Config.VerifySettings()
}
