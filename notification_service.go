package passport

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

type NotificationService struct {
	repository *NotificationRepository
	conf       *Config
	log        *zap.Logger
	mc         *MailClient
}

func NewNotificationService(repository *NotificationRepository, conf *Config, log *zap.Logger, mc *MailClient) *NotificationService {
	return &NotificationService{repository: repository, conf: conf, log: log, mc: mc}
}

func (s *NotificationService) Listener() {
	go func() {
		for {

			//TODO:for now only emails will not be sent
			notifications, err := s.repository.GetAllNotSent(context.Background())
			if err != nil {
				s.log.Error("could not get notifications")
				return
			}

			for _, notification := range notifications {
				email, found := notification.Params["email"]
				if !found {
					s.log.Error("email parameter not found")
					return
				}

				err = s.mc.Send(context.Background(), notification.Header, notification.Body, email)
				if err != nil {
					s.log.Error("could not send email")
					return
				}

				now := time.Now().UTC()

				err = s.repository.UpdateById(context.Background(), notification.ID, bson.M{"isSent": true, "sentOn": now})
				if err != nil {
					s.log.Error("could not update notification")
					return
				}
			}

			time.Sleep(5 * time.Second)
		}
	}()
}
