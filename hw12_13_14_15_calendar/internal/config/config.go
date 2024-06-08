package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

var options *Config

type Config struct {
	Logger     *LoggerConfig
	Database   *DatabaseConfig
	HTTPServer *ServerConfig
	GRPCServer *ServerConfig
	RabbitMQ   *RabbitMQConfig
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

type ConnectionConfig struct {
	Login    string
	Password string
	Host     string
	Port     string
}

type ConsumeConfig struct {
	Queue     string
	Consumer  string
	AutoAck   bool
	Exclusive bool
	NoLocal   bool
	NoWait    bool
	Interval  time.Duration
}

type PublishConfig struct {
	Exchange    string
	Key         string
	ContentType string
	Mandatory   bool
	Immediate   bool
}

type QueueConfig struct {
	Name       string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
}

type RabbitMQConfig struct {
	Connection *ConnectionConfig
	Publish    *PublishConfig
	Queue      *QueueConfig
	Consume    *ConsumeConfig
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
		Logger: &LoggerConfig{
			Level: viper.GetString("logger.level"),
		},
		Database: &DatabaseConfig{
			Prefix:       viper.GetString("database.Prefix"),
			DatabaseName: viper.GetString("database.DatabaseName"),
			Host:         viper.GetString("database.Host"),
			Port:         viper.GetString("database.Port"),
			UserName:     viper.GetString("database.UserName"),
			Password:     viper.GetString("database.Password"),
		},
		HTTPServer: &ServerConfig{
			Host: viper.GetString("http_server.Host"),
			Port: viper.GetString("http_server.Port"),
		},
		GRPCServer: &ServerConfig{
			Host: viper.GetString("grpc_server.Host"),
			Port: viper.GetString("grpc_server.Port"),
		},
		RabbitMQ: &RabbitMQConfig{
			Connection: &ConnectionConfig{
				Login:    viper.GetString("connection.login"),
				Password: viper.GetString("connection.password"),
				Host:     viper.GetString("connection.host"),
				Port:     viper.GetString("connection.port"),
			},
			Publish: &PublishConfig{
				Exchange:    viper.GetString("publish.exchange"),
				Key:         viper.GetString("publish.key"),
				ContentType: viper.GetString("publish.contentType"),
				Mandatory:   viper.GetBool("publish.mandatory"),
				Immediate:   viper.GetBool("publish.immediate"),
			},
			Queue: &QueueConfig{
				Name:       viper.GetString("queue.name"),
				Durable:    viper.GetBool("queue.durable"),
				AutoDelete: viper.GetBool("queue.autoDelete"),
				Exclusive:  viper.GetBool("queue.exclusive"),
				NoWait:     viper.GetBool("queue.noWait"),
			},
			Consume: &ConsumeConfig{
				Queue:     viper.GetString("consume.queue"),
				Consumer:  viper.GetString("consume.consumer"),
				AutoAck:   viper.GetBool("consume.autoAck"),
				Exclusive: viper.GetBool("consume.exclusive"),
				NoLocal:   viper.GetBool("consume.noLocal"),
				NoWait:    viper.GetBool("consume.noWait"),
				Interval:  viper.GetDuration("consume.interval"),
			},
		},
	}, nil
}

func Get() *Config {
	return options
}

func (config *LoggerConfig) GetLevel() string {
	return config.Level
}

func LoadConfig(configPath string) error {
	config, err := New(configPath)
	if err != nil {
		return err
	}
	options = &config
	return nil
}

func (s *ServerConfig) GetPort() string {
	return s.Port
}

func (s *ServerConfig) GetHost() string {
	return s.Host
}
