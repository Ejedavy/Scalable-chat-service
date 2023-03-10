package utils

import (
	"github.com/spf13/viper"
)

type Config struct {
	UsersList       string `mapstructure:"USERS_LIST"`
	UserChannelsFmt string `mapstructure:"USER_CHANNELS_FORMAT"`
	GeneralChannels string `mapstructure:"GENERAL_CHANNELS"`
	RedisAddress    string `mapstructure:"REDIS_ADDRESS"`
	ServerAddress   string `mapstructure:"SERVER_ADDRESS"`
}

func NewConfig() (Config, error) {
	config := Config{}
	viper.SetConfigName("app") // name of config file (without extension)
	viper.SetConfigType("env") // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()

	if err != nil {
		return config, err
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return config, err
	}
	return config, nil
}
