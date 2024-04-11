package util

import "github.com/spf13/viper"

// Config stores all configuration of the application
// The values are read by viper from a config file or environment variables.

type Config struct {
	DBDriver               string `mapstructure:"DB_DRIVER"`
	DBSource               string `mapstructure:"DB_SOURCE"`
	ServerAddress          string `mapstructure:"SERVER_ADDRESS"`
	Auth0Audience          string `mapstructure:"AUTH0_AUDIENCE"`
	Auth0Domain            string `mapstructure:"AUTH0_DOMAIN"`
	SigningKey             string `mapstructure:"SIGNING_KEY"`
	TwilioAccountSID       string `mapstructure:"TWILIO_ACCOUNT_SID"`
	TwilioAuthToken        string `mapstructure:"TWILIO_AUTH_TOKEN"`
	VerifyMobileServiceSID string `mapstructure:"VERIFY_MOBILE_SERVICE_SID"`
	VerifyEmailServiceSID  string `mapstructure:"VERIFY_EMAIL_SERVICE_SID"`
}

// LoadConfig reads configuration from file or environment variables
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
