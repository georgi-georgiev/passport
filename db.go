package passport

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoClient(conf *Config) *mongo.Client {
	credentials := options.Credential{Username: conf.Mongo.Username, Password: conf.Mongo.Password}
	clientOptions := options.Client().ApplyURI(conf.Mongo.Url).SetAuth(credentials)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		panic(err)
	}
	return client
}
