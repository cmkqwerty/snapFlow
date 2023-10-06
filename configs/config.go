package configs

import (
	"github.com/spf13/viper"
)

type Config struct {
	Environment   string `mapstructure:"ENVIRONMENT"`
	DBHost        string `mapstructure:"DB_HOST"`
	DBPort        string `mapstructure:"DB_PORT"`
	DBUser        string `mapstructure:"DB_USER"`
	DBPassword    string `mapstructure:"DB_PASSWORD"`
	DBName        string `mapstructure:"DB_NAME"`
	DBSSLMode     string `mapstructure:"DB_SSLMODE"`
	CSRFKey       string `mapstructure:"CSRF_KEY"`
	CSRFSecure    string `mapstructure:"CSRF_SECURE"`
	SMTPHost      string `mapstructure:"SMTP_HOST"`
	SMTPPort      string `mapstructure:"SMTP_PORT"`
	SMTPUsername  string `mapstructure:"SMTP_USERNAME"`
	SMTPPassword  string `mapstructure:"SMTP_PASSWORD"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
}

/*
func init() {

}

func ReadKey(key string) string {
	return viper.GetString(key)
}

*/

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
