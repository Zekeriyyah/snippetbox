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
	//Retrieve the id and hashed passwd for the provided email, if not exists return ErrInvalidCredentials
	var id int
	var hashedPasswd []byte

	stmt := `SELECT id, hashed_password FROM users where email=?`

	err := m.DB.QueryRow(stmt, email).Scan(&id, &hashedPasswd)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	//Check if the provided password and match the hashedPassword from the database if not
	//return ErrInvalidCredential error
	err = bcrypt.CompareHashAndPassword(hashedPasswd, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	//If no any error at all, return the id associated with the email
	return id, nil
}

// Implement Exists() to check for existence of user in the db, if user exists return true
func (m *UserModel) Exists(id int) (bool, error) {
	var exists bool

	stmt := "SELECT EXISTS(SELECT true FROM users WHERE id=?)"
	err := m.DB.QueryRow(stmt, id).Scan(&exists)

	return exists, err
}
