package config

type Config struct {
	OAuth2 GiteaOauth2Config  `yaml:"oauth2"`
	Remote RemoteConfigCenter `yaml:"remote"`
	DB     DBConfig           `yaml:"db"`
	Redis  RedisConfig        `yaml:"redis"`
}

type RemoteConfigCenter struct {
	Provider string `yaml:"provider"`
	EndPoint string `yaml:"endpoint"`
	Path     string `yaml:"path"`
}

type DBConfig struct {
	DSN string `yaml:"dsn"`
}

type RedisConfig struct {
	Addr string `yaml:"addr"`
}

type GiteaOauth2Config struct {
	BaseURL      string `yaml:"baseURL"`
	ClientID     string `yaml:"clientID"`
	ClientSecret string `yaml:"clientSecret"`
}
