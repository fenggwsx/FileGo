package conf

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type ServerConfig struct {
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	ShutdownTimeout   time.Duration
	MaxHeaderBytes    int
}

type Config struct {
	Development bool
	Port        uint16
	Server      ServerConfig
}

func NewConfig(flagset *pflag.FlagSet) (*Config, error) {
	config := &Config{}

	// Bind to the flagset containing global configs.
	viper.BindPFlags(flagset)

	// Set default values.
	viper.SetDefault("development", false)
	viper.SetDefault("port", 8080)
	viper.SetDefault("server.shutdownTimeout", "5s")

	// Bind environment variables.
	viper.MustBindEnv("development")
	viper.MustBindEnv("port")
	viper.MustBindEnv("server.readTimeout")
	viper.MustBindEnv("server.readHeaderTimeout")
	viper.MustBindEnv("server.writeTimeout")
	viper.MustBindEnv("server.idleTimeout")
	viper.MustBindEnv("server.shutdownTimeout")
	viper.MustBindEnv("server.maxHeaderBytes")

	// Set the location of the config file.
	viper.AddConfigPath(".")
	viper.SetConfigName("settings")

	// Read config data from file.
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	// Unmarshal the config data from viper.
	if err := viper.Unmarshal(config); err != nil {
		return nil, err
	}

	// Set gin mode according to global settings.
	if config.Development {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	return config, nil
}
