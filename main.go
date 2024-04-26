package main

import (
	"booking/config"
	"booking/internal/app"
	"booking/internal/cache"
	"booking/internal/environment"
	"booking/internal/logs"

	"github.com/gofiber/storage/redis/v3"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	err := environment.InitializeEnvironmentVariable()
	if err != nil {
		panic(err)
	}
	redis := redis.New(config.RedisConfig)
	cache := cache.New(redis)
	db, err := gorm.Open(postgres.New(config.PostgresConfig()), &gorm.Config{})
	if err != nil {
		logs.Error(err)
		panic("Cannot connect to the database")
	}

	app := app.GetApp(cache, db)
	port := viper.GetString("app.port")
	logs.Info("start at port " + port)
	app.Listen(":" + port)
}
