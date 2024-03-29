package facade

import (
	"context"
	"time"

	"github.com/georgi-georgiev/passport/notifications"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type Message struct {
	Topic     string
	Header    string
	Body      string
	Params    map[string]string
	Meta      []string
	Timestamp int64
}

type NotificationFacade struct {
	repository *notifications.NotificationRepository
	log        *zap.Logger
}

func NewNotificationFacade(repository *notifications.NotificationRepository, log *zap.Logger) *NotificationFacade {
	return &NotificationFacade{repository: repository, log: log}
}

func (f *NotificationFacade) Publish(ctx context.Context, flow string, message Message, userId string) {
	now := time.Now().UTC()
	notification := &notifications.Notification{
		ID:        primitive.NewObjectID(),
		CreatedOn: now,
		Flow:      flow,
		Topic:     message.Topic,
		Header:    message.Header,
		Body:      message.Body,
		Params:    message.Params,
		Meta:      message.Meta,
		UserID:    userId,
		IsSent:    false,
		IsRead:    false,
	}

	ts := time.Unix(message.Timestamp, 0)

	if flow == "email" {
		notification.IsRead = true
	} else {
		notification.IsSent = true
		notification.SentOn = &ts
	}

	_, err := f.repository.Create(ctx, notification)
	if err != nil {
		f.log.With(zap.Error(err)).Error("could not create notification")
		return
	}
}
