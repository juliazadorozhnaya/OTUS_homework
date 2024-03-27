package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Logger   LoggerConfig
	Database DatabaseConfig
	Server   ServerConfig
}

type LoggerConfig struct {
	Level string
}

type DatabaseConfig struct {
	Prefix       string
	DatabaseName string
	Host         string
	Port         string
	UserName     string
	Password     string
}

type ServerConfig struct {
	Host string
	Port string
}

func New(configPath string) (Config, error) {
	viper.SetConfigFile(configPath)

	if err := viper.ReadInConfig(); err != nil {
		return Config{}, fmt.Errorf("fatal error config file: %w", err)
	}

	return Config{
		Logger: LoggerConfig{
			Level: viper.GetString("logger.level"),
		},
		Database: DatabaseConfig{
			Prefix:       viper.GetString("database.Prefix"),
			DatabaseName: viper.GetString("database.DatabaseName"),
			Host:         viper.GetString("database.Host"),
			Port:         viper.GetString("database.Port"),
			UserName:     viper.GetString("database.UserName"),
			Password:     viper.GetString("database.Password"),
		},
		Server: ServerConfig{
			Host: viper.GetString("server.Host"),
			Port: viper.GetString("server.Port"),
		},
	}, nil
}
