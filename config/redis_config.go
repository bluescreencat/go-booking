package config

import (
	"github.com/gofiber/storage/redis/v3"
	"github.com/spf13/viper"
)

var RedisConfig = redis.Config{
	Host:     viper.GetString("redis.host"),
	Port:     viper.GetInt("redis.port"),
	Username: viper.GetString("redis.username"),
	Password: viper.GetString("redis.password"),
}
