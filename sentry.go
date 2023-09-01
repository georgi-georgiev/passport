package passport

import (
	"github.com/getsentry/sentry-go"
)

func NewSentry(conf *Config) *sentry.Client {
	sentryClient, err := sentry.NewClient(sentry.ClientOptions{
		Dsn:           conf.Sentry.DSN,
		EnableTracing: true,
	})

	if err != nil {
		panic(err)
	}

	return sentryClient
}
