package internal

import (
	"crypto/sha1"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"math/rand"
)

func InitializeAndGetEnv(key string) (string, error) {
	viper.SetConfigFile("$HOME/hamed/Projects/Own/Authentication_Microservice/.env")
	err := viper.ReadInConfig()
	if err != nil {
		return "", errors.Wrap(err, "viper can't read the config file")
	}
	value, ok := viper.Get(key).(string)
	if !ok {
		return "", fmt.Errorf("type assertion failed for the key: %s", key)
	}
	return value, nil
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

func GenerateCustomKey(email, hashToken string) string {
	h := sha1.New()
	h.Write([]byte(email + hashToken))
	byteSlice := h.Sum(nil)
	return string(byteSlice)
}
