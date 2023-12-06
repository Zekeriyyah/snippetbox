package models

import "errors"

var (
	ErrNoRecord = errors.New("models: no matching record found")

	//Add a ErrInvalidCredentials error
	ErrInvalidCredentials = errors.New("models: invalid credentials")

	//Add an ErrDuplicateEmail error
	ErrDuplicateEmail = errors.New("models: duplicate email")
)
