package models

import "errors"

var (
	// activities
	ErrActivityExists   = errors.New("activity already exists")
	ErrActivityNotFound = errors.New("activity not found")
	ErrForbidden        = errors.New("forbidden")

	// users
	ErrUserExists   = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")
)
