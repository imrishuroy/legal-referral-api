package util

import "github.com/spf13/viper"

// Config stores all configuration of the application
// The values are read by viper from a config file or environment variables.

type Config struct {
	Env                    string `mapstructure:"ENV"`
	DBDriver               string `mapstructure:"DB_DRIVER"`
	DBSource               string `mapstructure:"DB_SOURCE"`
	ServerAddress          string `mapstructure:"SERVER_ADDRESS"`
	FirebaseAuthKey        string `mapstructure:"FIREBASE_AUTH_KEY"`
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
	AWSBucketName          string `mapstructure:"AWS_BUCKET_NAME"`
	CloudFrontURL          string `mapstructure:"CLOUDFRONT_URL"`
	LinkedinClientID       string `mapstructure:"LINKEDIN_CLIENT_ID"`
	LinkedinClientSecret   string `mapstructure:"LINKEDIN_CLIENT_SECRET"`
	BootStrapServers       string `mapstructure:"BOOTSTRAP_SERVERS"`
	SecurityProtocol       string `mapstructure:"SECURITY_PROTOCOL"`
	SASLMechanism          string `mapstructure:"SASL_MECHANISM"`
	SASLUsername           string `mapstructure:"SASL_USERNAME"`
	SASLPassword           string `mapstructure:"SASL_PASSWORD"`
	Topic                  string `mapstructure:"TOPIC"`
	RedisHost              string `mapstructure:"REDIS_HOST"`
	RedisPort              string `mapstructure:"REDIS_PORT"`
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
