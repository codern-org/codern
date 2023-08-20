package domain

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
	InfluxDb ConfigInfluxDb `yaml:"influxdb" validate:"required"`
	MySql    ConfigMySql    `yaml:"mysql" validate:"required"`
	Fiber    ConfigFiber    `yaml:"fiber" validate:"required"`
	Frontend ConfigFrontend `yaml:"frontend" validate:"required"`
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

type ConfigFiber struct {
	Address string `yaml:"address" validate:"required"`
}

type ConfigFrontend struct {
	BaseUrl string             `yaml:"baseUrl" vallidate:"url,required"`
	Path    ConfigFrontendPath `yaml:"path"`
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
