package handlers

import "github.com/coke-day/functions/users/model"

// UserRegistration model for input request
type UserRegistration struct {
	Name         string `json:"name" validate:"required,min=2,max=100"`
	Email        string `json:"email" validate:"required,email,emailCustom"`
	Password     string `json:"password" validate:"required,min=2,max=100"` // TODO: Add more requirements to password
	HashPassword []byte `json:"-"`
}

// UserLogin model for input request
type UserLogin struct {
	Email        string `json:"email"  validate:"required,email,emailCustom"`
	Password     string `json:"password" validate:"required,min=2,max=100"`
	HashPassword []byte `json:"-"`
}

func mapRegisterToUser(u UserRegistration, salt string) model.User {
	return model.User{
		Email:        u.Email,
		HashPassword: hashPassword(u.Password, salt),
		Name:         u.Name,
	}
}

func mapLoginToUser(u UserLogin, salt string) model.User {
	return model.User{
		Email:        u.Email,
		HashPassword: hashPassword(u.Password, salt),
	}
}
