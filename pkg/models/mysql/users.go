package mysql

import (
	"database/sql"

	"jackson.software/snippetbox/pkg/models"
)

type UserRepository struct {
	DB *sql.DB
}

// Insert inserts a new user into the database. If there is already a user
// by the give data, an error will be returned.
func (m *UserRepository) Insert(name, email, password string) error {
	return nil
}

// Authenticate authenticates a user. If the authentication was not successfull,
// an error will be returned.
func (m *UserRepository) Authenticate(email, password string) (int, error) {
	return 0, nil
}

// Get gets a user found by the given ID. If no user could be found by
// the given ID, an error will be returned.
func (m *UserRepository) Get(id int) (*models.User, error) {
	return nil, nil
}
