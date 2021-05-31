package model

import (
	"github.com/coke-day/pkg/datastore"
)

// NewClientRepository instance
func NewClientRepository(ds datastore.Datastore) *ClientRepository {
	return &ClientRepository{datastore: ds}
}

// ClientRepository stores and fetches items
type ClientRepository struct {
	datastore datastore.Datastore
}

// Login a client
func (r *ClientRepository) Login(userLogin *UserLogin) (*UserPersistence, error) {
	user := &UserPersistence{}
	if err := r.datastore.Get(userLogin.Email, getSecondaryKey(), &user); err != nil {
		return nil, err
	}

	return user, nil
}

// Register a new client
func (r *ClientRepository) Register(userRegister *UserRegistration) error {
	user := userRegister.toUserPersistence()
	return r.datastore.Store(user)
}
