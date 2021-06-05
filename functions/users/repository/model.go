package repository

import (
	"github.com/coke-day/functions/users/model"
)

// instanceItem describes User dynamodb representation
type instanceItem struct {
	PK           string `dynamodbav:"pk"`
	SK           string `dynamodbav:"sk"`
	Email        string `dynamodbav:"email"`
	HashPassword []byte `dynamodbav:"hash_password"`
	Name         string `dynamodbav:"name"`
}

func getSecondaryKey() string {
	return "users"
}

func buildInstanceItem(u model.User) instanceItem {
	return instanceItem{
		PK:           u.Email,
		SK:           getSecondaryKey(),
		Email:        u.Email,
		HashPassword: u.HashPassword,
		Name:         u.Name,
	}
}

func mapInstanceItemToUser(instanceItem instanceItem) *model.User {
	return &model.User{
		Email:        instanceItem.Email,
		Name:         instanceItem.Name,
		HashPassword: instanceItem.HashPassword,
	}
}
