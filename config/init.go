package config

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

var Cfg Config

func InitConfig(defaultRelConfigPath string) {
	cfile := pflag.String("config", defaultRelConfigPath, "config file path")
	pflag.Parse()
	viper.SetConfigFile(*cfile)
	viper.WatchConfig()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	// initViperRemote()

	err = viper.Unmarshal(&Cfg)
	if err != nil {
		panic(err)
	}
}

// func initViperRemote() {
// 	if val := viper.Get("remote"); val == nil {
// 		// No remote config center
// 		slog.Warn("remote config center is not set")
// 		return
// 	}
// 	err := viper.AddRemoteProvider(
// 		viper.GetString("remote.provider"),
// 		viper.GetString("remote.endpoint"),
// 		viper.GetString("remote.path"),
// 	)
// 	if err != nil {
// 		panic(err)
// 	}
// 	viper.SetConfigType("yaml")
// 	err = viper.ReadRemoteConfig()
// 	if err != nil {
// 		panic(err)
// 	}
// }
