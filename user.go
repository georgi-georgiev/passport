package passport

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID                primitive.ObjectID   `bson:"_id"`
	CreatedOn         time.Time            `bson:"createdOn"`
	UpdatedOn         *time.Time           `bson:"updatedOn,omitempty"`
	IsActive          bool                 `bson:"isActive"`
	IsVerified        bool                 `bson:"isVerified"`
	VerificationToken string               `bson:"verificationToken,omitempty"`
	Username          string               `bson:"username"`
	Email             string               `bson:"email"`
	Password          string               `bson:"password"`
	RecoveryCode      string               `bson:"recoveryCode,omitempty"`
	ResettingCode     string               `bson:"resettingCode,omitempty"`
	Role              primitive.ObjectID   `bson:"role,omitempty"`
	Rights            []primitive.ObjectID `bson:"rights,omitempty"`
}

func NewUser(verificationToken string, username string, email string, passwordHash string, role *Role, rights []*Right) *User {

	rightsIds := make([]primitive.ObjectID, 0)
	for _, right := range rights {
		rightsIds = append(rightsIds, right.ID)
	}

	return &User{
		ID:                primitive.NewObjectID(),
		CreatedOn:         time.Now().UTC(),
		IsActive:          true,
		IsVerified:        false,
		VerificationToken: verificationToken,
		Username:          username,
		Email:             email,
		Password:          passwordHash,
		Role:              role.ID,
		Rights:            rightsIds,
	}
}
