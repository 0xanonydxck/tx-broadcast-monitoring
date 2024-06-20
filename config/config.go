package config

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var config Config

func Init() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Panic().Err(err).Msg("config::Init(): failed to read config")
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Panic().Err(err).Msg("config::Init(): failed to unmarshal config")
	}
}

func Get() Config {
	return config
}
