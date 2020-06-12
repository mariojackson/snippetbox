package mysql

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"jackson.software/snippetbox/pkg/models"
)

type UserRepository struct {
	DB *sql.DB
}

// Insert inserts a new user into the database. If there is already a user
// by the give data, an error will be returned.
func (m *UserRepository) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	var duplicateEntryCode uint16 = 1062

	stmt := `INSERT INTO users (name, email, hashed_password, created)
	VALUES (?, ?, ?, UTC_TIMESTAMP())`

	_, err = m.DB.Exec(stmt, name, email, string(hashedPassword))
	if err != nil {
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == duplicateEntryCode && strings.Contains(mySQLError.Message, "users_uc_email") {
				return models.ErrDuplicateEmail
			}
		}
		return err
	}

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
