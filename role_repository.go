package passport

import (
	"context"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RoleRepository struct {
	*MongoRepository
}

func NewRoleRepository(client *mongo.Client, conf *Config) *RoleRepository {
	repository := NewMongoRepository(client, conf.Mongo.Dbname, "roles")

	nameIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "name", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	_, err := repository.Collection.Indexes().CreateMany(context.TODO(), []mongo.IndexModel{nameIndex})
	if err != nil {
		panic(err)
	}

	adminRole := &Role{
		Name: "admin",
	}

	filter := bson.M{"name": adminRole.Name}
	update := bson.M{"$setOnInsert": adminRole}

	opts := options.Update().SetUpsert(true)
	_, err = repository.Collection.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		panic(err)
	}

	return &RoleRepository{repository}
}

func (r *RoleRepository) GetByName(ctx context.Context, name string) (*Role, error) {
	role := &Role{}

	err := r.Collection.FindOne(ctx, bson.M{
		"name": strings.ToLower(name),
	}).Decode(role)

	if err != nil {
		return nil, err
	}

	return role, nil
}

func (r *RoleRepository) GetAll(ctx context.Context) ([]*Role, error) {
	roles := []*Role{}

	cur, err := r.Collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	defer cur.Close(ctx)

	for cur.Next(ctx) {
		r := &Role{}

		err = cur.Decode(r)

		if err != nil {
			return nil, err
		}

		roles = append(roles, r)
	}

	return roles, nil
}

func (r *RoleRepository) Update(ctx context.Context, role *Role) error {

	updateBody := bson.M{}

	updateBody["name"] = role.Name

	return r.UpdateById(ctx, role.ID, updateBody)
}
