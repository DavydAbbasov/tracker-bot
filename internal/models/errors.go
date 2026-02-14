package models

import "errors"

var (
	// Activity domain errors.
	ErrActivityExists   = errors.New("activity already exists")
	ErrActivityNotFound = errors.New("activity not found")
	ErrForbidden        = errors.New("forbidden")

	// User domain errors.
	ErrUserExists   = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")
)
