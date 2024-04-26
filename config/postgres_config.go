package config

import (
	"fmt"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
)

func PostgresConfig() postgres.Config {
	dsn := fmt.Sprintf(
		`host=%v user=%v password=%v dbname=%v port=%v sslmode=%v TimeZone=%v`,
		viper.GetString("db.host"),
		viper.GetString("db.username"),
		viper.GetString("db.password"),
		viper.GetString("db.name"),
		viper.GetInt("db.port"),
		viper.GetString("db.sslmode"),
		viper.Get("db.timezone"),
	)

	return postgres.Config{
		DSN: dsn,
	}
}
