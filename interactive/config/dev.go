// unused. test for build tag
//go:build dev

package config

var Config = config{
	DB: DBConfig{
		DSN: "root:root@tcp(127.0.0.1:13316)/wetravel?charset=utf8mb4&parseTime=True&loc=Local",
	},
	Redis: RedisConfig{
		Addr: "localhost:6379",
	},
}
