package passport

import (
	"context"

	"github.com/rotisserie/eris"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepository struct {
	Client     *mongo.Client
	Database   *mongo.Database
	Collection *mongo.Collection
}

func NewMongoRepository(client *mongo.Client, dbname string, collectionName string) *MongoRepository {
	database := client.Database(dbname)
	return &MongoRepository{Client: client, Database: database, Collection: database.Collection(collectionName)}
}

func (r *MongoRepository) Create(ctx context.Context, v interface{}) (primitive.ObjectID, error) {
	result, err := r.Collection.InsertOne(ctx, v)
	if err != nil {
		return primitive.NilObjectID, eris.Wrap(err, "could not insert")
	}

	newID := result.InsertedID.(primitive.ObjectID)

	return newID, nil
}

func (r *MongoRepository) GetById(ctx context.Context, id primitive.ObjectID, v interface{}) (bool, error) {

	err := r.Collection.FindOne(ctx, bson.M{
		"_id": id,
	}).Decode(v)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (r *MongoRepository) UpdateById(ctx context.Context, id primitive.ObjectID, updateBody bson.M) error {
	_, err := r.Collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": updateBody})
	return err
}

func (r *MongoRepository) GetStringFieldForId(ctx context.Context, id primitive.ObjectID, field string) (string, error) {
	var result bson.M

	opt := options.FindOne().SetProjection(bson.M{"_id": 0, field: 1})
	err := r.Collection.FindOne(ctx, bson.M{"_id": id}, opt).Decode(&result)

	if err != nil {
		return "", err
	}

	return result[field].(string), nil
}

func (r *MongoRepository) DeleteById(ctx context.Context, id primitive.ObjectID) (bool, error) {
	dr, err := r.Collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return false, err
	}

	return dr.DeletedCount > 0, nil
}

func (r *MongoRepository) SetFieldAndWipeOtherForId(ctx context.Context, id primitive.ObjectID, fieldToSet string, value string, fieldToWipe string) error {
	updateBody := bson.M{fieldToSet: value, fieldToWipe: ""}

	_, err := r.Collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": updateBody})
	return err
}
