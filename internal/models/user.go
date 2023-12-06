package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type UserModel struct {
	DB *sql.DB
}

// Implement Insert to add user to the model
func (m *UserModel) Insert(name, email, password string) error {
	//Create bcrypt hash of the password
	hashedPasswd, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	// Query DB and execute the stmt to store name, email, hashedPasswd and date created into table users
	stmt := `INSERT INTO users (name, email, hashed_password, created)
	VALUES(?,?,?, UTC_TIMESTAMP())`

	_, err = m.DB.Exec(stmt, name, email, string(hashedPasswd))
	if err != nil {
		//Check if the error is mySQLError
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return ErrDuplicateEmail
			}
		}
		return err
	}

	return nil
}

// Implement Authenticate to verify whether user with provided email and password exists
func (m *UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

// Implement Exists() to check for existence of user in the db, if user exists return true
func (m *UserModel) Exists(id int) (bool, error) {
	return false, nil
}
