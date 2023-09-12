package notifications

import (
	"context"

	"github.com/georgi-georgiev/passport"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type NotificationRepository struct {
	*passport.MongoRepository
}

func NewNotificationRepository(client *mongo.Client, conf *passport.Config) *NotificationRepository {
	repository := passport.NewMongoRepository(client, conf.Mongo.Dbname, "notifications")
	return &NotificationRepository{repository}
}

func (r *NotificationRepository) GetAllNotSent(ctx context.Context) ([]*Notification, error) {
	notifications := []*Notification{}
	filter := bson.M{"isSent": false}
	cur, err := r.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	defer cur.Close(ctx)

	for cur.Next(ctx) {
		notification := &Notification{}

		err = cur.Decode(notification)

		if err != nil {
			return nil, err
		}

		notifications = append(notifications, notification)
	}

	return notifications, nil
}

func (r *NotificationRepository) GetAllNotRead(ctx context.Context) ([]*Notification, error) {
	notifications := []*Notification{}
	filter := bson.M{"isRead": false}
	cur, err := r.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	defer cur.Close(ctx)

	for cur.Next(ctx) {
		notification := &Notification{}

		err = cur.Decode(notification)

		if err != nil {
			return nil, err
		}

		notifications = append(notifications, notification)
	}

	return notifications, nil
}
