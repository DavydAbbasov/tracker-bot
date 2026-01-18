package errors

import "errors"

var (
	ErrActivityExists   = errors.New("activity already exists")
	ErrActivityNotFound = errors.New("activity not found")
	ErrForbidden        = errors.New("forbidden")
)
