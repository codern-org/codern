package domain

type Config struct {
	Metadata ConfigMetadata `yaml:"metadata"`
	Client   ConfigClient   `yaml:"client"`
	Google   ConfigGoogle   `yaml:"google"`
	Auth     ConfigAuth     `yaml:"auth"`
}

type ConfigMetadata struct {
	Name     string `yaml:"name"`
	Version  string `yaml:"version"`
	LogLevel string `yaml:"logLevel"`
}

type ConfigClient struct {
	InfluxDb ConfigInfluxDb `yaml:"influxdb"`
	MySql    ConfigMySql    `yaml:"mysql"`
	Fiber    ConfigFiber    `yaml:"fiber"`
}

type ConfigInfluxDb struct {
	Url    string `yaml:"url"`
	Token  string `yaml:"token"`
	Org    string `yaml:"org"`
	Bucket string `yaml:"bucket"`
}

type ConfigMySql struct {
	Uri string `yaml:"uri"`
}

type ConfigFiber struct {
	Address string `yaml:"address"`
}

type ConfigGoogle struct {
	ClientId     string `yaml:"clientId"`
	ClientSecret string `yaml:"clientSecret"`
	RedirectUri  string `yaml:"redirectUri"`
}

type ConfigAuth struct {
	Session ConfigAuthSession `yaml:"session"`
}

type ConfigAuthSession struct {
	Prefix string `yaml:"prefix"`
	Secret string `yaml:"secret"`
	MaxAge int    `yaml:"maxAge"`
}
