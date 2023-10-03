package configs

import (
	"github.com/spf13/viper"
	"log"
)

func init() {
	viper.SetConfigFile("app.env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("can not load config:", err)
	}
}

func ReadKey(key string) string {
	return viper.GetString(key)
}
