package users

import (
	"context"
	"fmt"
	"strings"

	"github.com/georgi-georgiev/passport"
	"github.com/georgi-georgiev/passport/permissions"
	"github.com/rotisserie/eris"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepository struct {
	*passport.MongoRepository
}

func NewUserRepository(client *mongo.Client, conf *passport.Config, roleRepository *permissions.RoleRepository) *UserRepository {
	repository := passport.NewMongoRepository(client, conf.Mongo.Dbname, "users")

	usernameIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "username", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	emailIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	_, err := repository.Collection.Indexes().CreateMany(context.TODO(), []mongo.IndexModel{usernameIndex, emailIndex})
	if err != nil {
		panic(err)
	}

	hashedPassword, err := passport.Hash("admin")
	if err != nil {
		panic(err)
	}

	adminRole, err := roleRepository.GetByName(context.TODO(), "admin")
	if err != nil {
		panic(err)
	}

	adminUser := NewUser("", "admin", "admin@test.com", hashedPassword, adminRole, nil)

	adminUser.IsVerified = true

	filter := bson.M{"username": adminUser.Username}
	update := bson.M{"$setOnInsert": adminUser}

	opts := options.Update().SetUpsert(true)
	_, err = repository.Collection.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		panic(err)
	}

	return &UserRepository{repository}
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*User, error) {
	result := &User{}

	err := r.Collection.FindOne(ctx, bson.M{
		"username": strings.ToLower(username),
	}).Decode(result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	result := &User{}

	err := r.Collection.FindOne(ctx, bson.M{
		"email": strings.ToLower(email),
	}).Decode(result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return result, nil
}

func (r *UserRepository) Verify(ctx context.Context, token string) error {
	fmt.Println("token", token)
	filter := bson.M{"verificationToken": token}
	update := bson.M{"$set": bson.M{"isVerified": true, "verificationToken": ""}}
	ur, err := r.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if ur.ModifiedCount == 0 {
		return eris.New("user with token not found")
	}

	return nil
}

func (r *UserRepository) GetAll(ctx context.Context) ([]*User, error) {
	users := []*User{}

	cur, err := r.Collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	defer cur.Close(ctx)

	for cur.Next(ctx) {
		u := &User{}

		err = cur.Decode(u)

		if err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	return users, nil
}

func (r *UserRepository) Update(ctx context.Context, u *User) error {

	updateBody := bson.M{}
	updateBody["email"] = u.Email
	updateBody["username"] = u.Username

	return r.UpdateById(ctx, u.ID, updateBody)
}

func (r *UserRepository) SetRecoveryCode(ctx context.Context, id primitive.ObjectID, code string) error {
	return r.SetFieldAndWipeOtherForId(ctx, id, "recoveryCode", code, "resettingCode")
}

func (r *UserRepository) GetRecoveryCode(ctx context.Context, id primitive.ObjectID) (string, error) {
	return r.GetStringFieldForId(ctx, id, "recoveryCode")
}

func (r *UserRepository) SetResettingCode(ctx context.Context, id primitive.ObjectID, code string) error {
	return r.SetFieldAndWipeOtherForId(ctx, id, "resettingCode", code, "recoveryCode")
}

func (r *UserRepository) GetResettingCode(ctx context.Context, id primitive.ObjectID) (string, error) {
	return r.GetStringFieldForId(ctx, id, "resettingCode")
}

func (r *UserRepository) ResetPassword(ctx context.Context, id primitive.ObjectID, passwordHash string) error {
	return r.SetFieldAndWipeOtherForId(ctx, id, "password", passwordHash, "resettingCode")
}
