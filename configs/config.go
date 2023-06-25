// Package configs is described all configuration for this source
package configs

import (
	"github.com/spf13/viper"
)

type config struct {
	DBDNSAuth  string `mapstructure:"DB_DNS_AUTH"`
	DBDNSMovie string `mapstructure:"DB_DNS_MOVIE"`
}

var Configs *config

func InitConfigs() {
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	viper.WatchConfig()

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	err = viper.Unmarshal(&Configs)
	if err != nil {
		panic(err)
	}
}
