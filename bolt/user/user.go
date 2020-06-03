package user

import (
	"github.com/boltdb/bolt"
	helper_reset_password "github.com/portainer/helper-reset-password"
	"github.com/portainer/helper-reset-password/bolt/internal"
)

const (
	// BucketName represents the name of the bucket where this service stores data.
	BucketName = "users"
)

// Service represents a service for managing endpoint data.
type Service struct {
	db *bolt.DB
}

// NewService creates a new instance of a service.
func NewService(db *bolt.DB) *Service {
	return &Service{
		db: db,
	}
}

// User returns a user by ID
func (service *Service) User(ID int) (*helper_reset_password.User, error) {
	var user helper_reset_password.User
	identifier := internal.Itob(ID)

	err := internal.GetObject(service.db, BucketName, identifier, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdateUser saves a user.
func (service *Service) UpdateUser(ID int, user *helper_reset_password.User) error {
	identifier := internal.Itob(ID)
	return internal.UpdateObject(service.db, BucketName, identifier, user)
}
