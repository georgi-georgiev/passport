package notifications

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Notification struct {
	ID        primitive.ObjectID `bson:"_id"`
	CreatedOn time.Time          `bson:"createdOn"`
	Flow      string             `bson:"type"`
	UserID    string             `bson:"userId"`
	Topic     string             `bson:"topic"`
	Header    string             `bson:"header"`
	Body      string             `bson:"body"`
	Params    map[string]string  `bson:"params,omitempty"`
	Meta      []string           `bson:"meta,omitempty"`
	IsSent    bool               `bson:"isSent"`
	SentOn    *time.Time         `bson:"sentOn,omitempty"`
	IsRead    bool               `bson:"isRead"`
	ReadOn    *time.Time         `bson:"readOn,omitempty"`
}
