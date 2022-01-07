package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Configuration struct {
	Environment string
	Port        int64
}

func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath("./config")
	// read config from ENV
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	// read config
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func GetConfig() Configuration {
	initConfig()
	config := Configuration{}
	err := viper.Unmarshal(&config)
	if err != nil {
		panic(err)
	}
	return config
}