package repository

import (
	"github.com/coke-day/functions/users/model"
	"github.com/coke-day/pkg/datastore"
)

// NewUserRepository instance
func NewUserRepository(ds datastore.Datastore) *UserRepository {
	return &UserRepository{datastore: ds}
}

// UserRepository stores and fetches items
type UserRepository struct {
	datastore datastore.Datastore
}

// Login a client
func (r *UserRepository) Login(user model.User) (*model.User, error) {
	instanceUser := &instanceItem{}
	if err := r.datastore.Get(user.Email, getSecondaryKey(), &instanceUser); err != nil {
		return nil, err
	}

	return mapInstanceItemToUser(*instanceUser), nil
}

// Register a new client
func (r *UserRepository) Register(user model.User) error {
	return r.datastore.Store(buildInstanceItem(user))
}
