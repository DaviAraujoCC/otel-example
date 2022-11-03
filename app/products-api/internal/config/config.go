package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func New() (error) {
	viper.SetDefault("LOG_LEVEL", "debug")
	viper.SetDefault("HTTP_PORT", "8080")
	viper.SetDefault("MYSQL_HOST", "localhost")
	viper.SetDefault("MYSQL_USER", "root")
	viper.SetDefault("MYSQL_PASSWORD", "root")
	viper.SetDefault("OTEL_EXPORTER_ENDPOINT", "")

	viper.SetConfigType("env")

	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		logrus.Println("config: .env file not found")
	}

	viper.AutomaticEnv()


	return nil

}