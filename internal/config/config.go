package config

import (
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Metadata ConfigMetadata `yaml:"metadata" validate:"required"`
	Client   ConfigClient   `yaml:"client" validate:"required"`
	Google   ConfigGoogle   `yaml:"google" validate:"required"`
	Auth     ConfigAuth     `yaml:"auth" validate:"required"`
}

type ConfigMetadata struct {
	Name     string `yaml:"name" validate:"required"`
	Version  string `yaml:"version" validate:"required"`
	LogLevel string `yaml:"logLevel" validate:"required"`
}

type ConfigClient struct {
	InfluxDb  ConfigInfluxDb  `yaml:"influxdb" validate:"required"`
	MySql     ConfigMySql     `yaml:"mysql" validate:"required"`
	SeaweedFs ConfigSeaweedFs `yaml:"seaweedfs" validate:"required"`
	RabbitMq  ConfigRabbitMq  `yaml:"rabbitmq" validate:"required"`
	Fiber     ConfigFiber     `yaml:"fiber" validate:"required"`
	Frontend  ConfigFrontend  `yaml:"frontend" validate:"required"`
}

type ConfigInfluxDb struct {
	Url    string `yaml:"url" validate:"url,required"`
	Token  string `yaml:"token" validate:"required"`
	Org    string `yaml:"org" validate:"required"`
	Bucket string `yaml:"bucket" validate:"required"`
}

type ConfigMySql struct {
	Uri string `yaml:"uri" validate:"required"`
}

type ConfigSeaweedFs struct {
	MasterUrl string                  `yaml:"masterUrl" validate:"url,required"`
	FilerUrls ConfigSeaweedFSFilerUrl `yaml:"filerUrls" validate:"required"`
}

type ConfigSeaweedFSFilerUrl struct {
	Internal string `yaml:"internal" validate:"url,required"`
	External string `yaml:"external" validate:"url,required"`
}

type ConfigRabbitMq struct {
	Url string `yaml:"url" validate:"url,required"`
}

type ConfigFiber struct {
	Address string `yaml:"address" validate:"required"`
}

type ConfigFrontend struct {
	BaseUrl string             `yaml:"baseUrl" vallidate:"url,required"`
	Path    ConfigFrontendPath `yaml:"path" validate:"required"`
}

type ConfigFrontendPath struct {
	SignIn string `yaml:"signIn" validate:"required"`
}

type ConfigGoogle struct {
	ClientId     string `yaml:"clientId" validate:"required"`
	ClientSecret string `yaml:"clientSecret" validate:"required"`
	RedirectUri  string `yaml:"redirectUri" validate:"required"`
}

type ConfigAuth struct {
	Session ConfigAuthSession `yaml:"session" validate:"required"`
}

type ConfigAuthSession struct {
	Prefix string `yaml:"prefix" validate:"required"`
	Secret string `yaml:"secret" validate:"required"`
	MaxAge int    `yaml:"maxAge" validate:"number,required"`
}

func Load(path string) (*Config, error) {
	if err := validatePath(path); err != nil {
		return nil, err
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config *Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	validate := validator.New()
	if err := validate.Struct(config); err != nil {
		return nil, err
	}

	return config, nil
}

func validatePath(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	if info.IsDir() {
		return fmt.Errorf("%s is a directory, not a file", path)
	}

	return nil
}
