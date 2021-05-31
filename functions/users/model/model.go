package model

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

// UserPersistence describes dynamodb representation
type UserPersistence struct {
	PK           string `dynamodbav:"pk"`
	SK           string `dynamodbav:"sk"`
	Email        string `dynamodbav:"email"`
	HashPassword []byte `dynamodbav:"hash_password"`
	Name         string `dynamodbav:"name"`
}

func (u UserRegistration) toUserPersistence() UserPersistence {
	return UserPersistence{
		PK:           u.Email,
		SK:           getSecondaryKey(),
		Email:        u.Email,
		HashPassword: u.HashPassword,
		Name:         u.Name,
	}
}

func getSecondaryKey() string {
	return "users"
}
