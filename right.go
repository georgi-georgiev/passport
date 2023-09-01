package passport

import "go.mongodb.org/mongo-driver/bson/primitive"

type Right struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Name string             `bson:"name,omitempty"`
}
