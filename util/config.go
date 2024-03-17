package util

import "github.com/spf13/viper"

// Config stores all configuration of the application
// The values are read by viper from a config file or environemnt variables.

type Config struct {
	DBDriver      string `mapstructure:"DB_DRIVER"`
	DBSource      string `mapstructure:"DB_SOURCE"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
	Auth0Audience string `mapstructure:"AUTH0_AUDIENCE"`
	Auth0Domain   string `mapstructure:"AUTH0_DOMAIN"`
	SigningKey    string `mapstructure:"SIGNING_KEY"`
}

// LoadConfig reads configuration from file or environemnt variables
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env") // json, xml

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
