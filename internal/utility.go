package internal

import (
	"fmt"
	"github.com/spf13/viper"
)

func InitializeAndGetEnv(key string) (string, error) {
	viper.SetConfigFile("$HOME/hamed/Projects/Own/Authentication_Microservice/.env")
	err := viper.ReadInConfig()
	if err != nil {
		return "", err
	}
	value, ok := viper.Get(key).(string)
	if !ok {
		return "", fmt.Errorf("type assertion failed for the key: %s", key)
	}
	return value, nil
}
