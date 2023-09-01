package passport

import (
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Server  ServerConfiguration
	Mongo   MongoConfiguration
	Sentry  SentryConfiguration
	Mail    MailConfiguration
	App     AppConfiguration
	Swagger SwaggerConfiguration
}

type ServerConfiguration struct {
	RootURL         string
	Host            string
	Port            string
	ShutdownTimeout int
	Timeout         struct {
		Server time.Duration
		Read   time.Duration
		Write  time.Duration
		Idle   time.Duration
	}
}

type SwaggerConfiguration struct {
	Host string
}

type MongoConfiguration struct {
	Url      string
	Dbname   string
	Username string
	Password string
}

type SentryConfiguration struct {
	DSN string
}

type MailConfiguration struct {
	ApiKey      string
	PartnerKey  string
	SenderEmail string
}

type AppConfiguration struct {
	Name        string
	Version     string
	PrivKeyPath string
	PubKeyPath  string
}

func NewConfig() *Config {
	var config *Config

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	path := os.Getenv("CONFIG_PATH")

	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		panic(err)
	}

	return config
}
