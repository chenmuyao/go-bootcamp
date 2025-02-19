// unused. test for build tag
//go:build k8s

package config

var Config = config{
	DB: DBConfig{
		DSN: "root:root@tcp(wetravel-mysql:3308)/wetravel?charset=utf8mb4&parseTime=True&loc=Local",
	},
	Redis: RedisConfig{
		Addr: "wetravel-redis:6380",
	},
}
