package config

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var Cfg Config

func InitConfig(defaultRelConfigPath string) {
	cfile := pflag.String("config", defaultRelConfigPath, "config file path")
	pflag.Parse()
	viper.SetConfigFile(*cfile)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	err = viper.Unmarshal(&Cfg)
	if err != nil {
		panic(err)
	}
}
