package permissions

import (
	"context"
	"strings"

	"github.com/georgi-georgiev/passport"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RightRepository struct {
	*passport.MongoRepository
}

func NewRightRepository(client *mongo.Client, conf *passport.Config) *RightRepository {
	repository := passport.NewMongoRepository(client, conf.Mongo.Dbname, "rights")

	nameIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "name", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	_, err := repository.Collection.Indexes().CreateMany(context.TODO(), []mongo.IndexModel{nameIndex})
	if err != nil {
		panic(err)
	}

	return &RightRepository{repository}
}

func (r *RightRepository) GetByName(ctx context.Context, name string) (*Right, error) {
	right := &Right{}

	err := r.Collection.FindOne(ctx, bson.M{
		"name": strings.ToLower(name),
	}).Decode(right)

	if err != nil {
		return nil, err
	}

	return right, nil
}

func (r *RightRepository) GetAll(ctx context.Context) ([]*Right, error) {
	rights := []*Right{}

	cur, err := r.Collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	defer cur.Close(ctx)

	for cur.Next(ctx) {
		r := &Right{}

		err = cur.Decode(r)

		if err != nil {
			return nil, err
		}

		rights = append(rights, r)
	}

	return rights, nil
}

func (r *RightRepository) Update(ctx context.Context, right *Right) error {

	updateBody := bson.M{}
	updateBody["name"] = right.Name

	return r.UpdateById(ctx, right.ID, updateBody)
}
