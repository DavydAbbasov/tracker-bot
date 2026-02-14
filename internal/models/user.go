package models

import "time"

// User is a persisted user record from database.
type User struct {
	TgUserID    int64
	UserName    *string
	PhoneNumber *string
	Email       *string
	Language    *string
	TimeZone    *string
	CreatedAt   time.Time
}
