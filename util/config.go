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
	AWSRegion              string `mapstructure:"AWS_REGION"`
	AWSAccessKeyID         string `mapstructure:"AWS_ACCESS_KEY_ID"`
	AWSSecretKey           string `mapstructure:"AWS_SECRET_KEY"`
	AWSBucketPrefix        string `mapstructure:"AWS_BUCKET_PREFIX"`
	LinkedinClientID       string `mapstructure:"LINKEDIN_CLIENT_ID"`
	LinkedinClientSecret   string `mapstructure:"LINKEDIN_CLIENT_SECRET"`
	BootStrapServers       string `mapstructure:"BOOTSTRAP_SERVERS"`
	SecurityProtocol       string `mapstructure:"SECURITY_PROTOCOL"`
	SASLMechanism          string `mapstructure:"SASL_MECHANISM"`
	SASLUsername           string `mapstructure:"SASL_USERNAME"`
	SASLPassword           string `mapstructure:"SASL_PASSWORD"`
	Topic                  string `mapstructure:"TOPIC"`
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
