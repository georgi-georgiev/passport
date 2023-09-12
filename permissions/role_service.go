package permissions

import (
	"context"

	"github.com/georgi-georgiev/passport/payloads"
	"github.com/rotisserie/eris"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoleService struct {
	repository *RoleRepository
}

func NewRoleService(repository *RoleRepository) *RoleService {
	return &RoleService{repository: repository}
}

func (s *RoleService) CreateRole(ctx context.Context, payload payloads.CreateRolePayload) (*Role, error) {

	role := &Role{
		Name: payload.Name,
	}

	_, err := s.repository.Create(ctx, role)
	if err != nil {
		return nil, err
	}

	return role, nil
}

func (s *RoleService) GetRoles(ctx context.Context) ([]*Role, error) {
	roles, err := s.repository.GetAll(ctx)
	if err != nil {
		return nil, eris.Wrap(err, "could no get all roles")
	}

	return roles, nil
}

func (s *RoleService) GetByName(ctx context.Context, name string) (*Role, error) {
	return s.repository.GetByName(ctx, name)
}

func (s *RoleService) GetById(ctx context.Context, id primitive.ObjectID) (*Role, error) {
	role := &Role{}
	_, err := s.repository.GetById(ctx, id, role)
	if err != nil {
		return nil, err
	}

	return role, nil
}

func (s *RoleService) UpdateRole(ctx context.Context, id primitive.ObjectID, name string) (*Role, error) {

	r := &Role{}
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
