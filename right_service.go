package passport

import (
	"context"

	"github.com/rotisserie/eris"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type RightService struct {
	repository *RightRepository
	conf       *Config
	log        *zap.Logger
}

func NewRightService(repository *RightRepository, conf *Config, log *zap.Logger) *RightService {
	return &RightService{repository: repository, conf: conf, log: log}
}

func (s *RightService) CreateRight(ctx context.Context, payload CreateRightPayload) (*Right, error) {

	right := &Right{
		Name: payload.Name,
	}

	_, err := s.repository.Create(ctx, right)
	if err != nil {
		return nil, err
	}

	return right, nil
}

func (s *RightService) GetByName(ctx context.Context, name string) (*Right, error) {
	return s.repository.GetByName(ctx, name)
}

func (s *RightService) GetManyByIds(ctx context.Context, ids []primitive.ObjectID) ([]Right, error) {

	filter := bson.M{"_id": bson.M{"$in": ids}}

	cursor, err := s.repository.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var rights []Right

	for cursor.Next(ctx) {
		var right Right
		if err := cursor.Decode(&right); err != nil {
			return nil, err
		}
		rights = append(rights, right)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return rights, nil
}

func (s *RightService) GetRights(ctx context.Context) ([]*Right, error) {
	rights, err := s.repository.GetAll(ctx)
	if err != nil {
		return nil, eris.Wrap(err, "could no get all rights")
	}

	return rights, nil
}

func (s *RightService) UpdateRight(ctx context.Context, id primitive.ObjectID, name string) (*Right, error) {

	r := &Right{}
	found, err := s.repository.GetById(ctx, id, r)
	if err != nil {
		return nil, err
	}

	if !found {
		return nil, nil
	}

	if r.Name != name {
		r.Name = name
	}

	err = s.repository.Update(ctx, r)
	if err != nil {
		return nil, err
	}

	return r, nil
}
