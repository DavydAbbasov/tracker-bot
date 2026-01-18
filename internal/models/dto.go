package models

import "time"

// для входящих
type UserInput struct {
	TgUserID    int64
	UserName    *string
	PhoneNumber *string
	Email       *string
	Language    *string
	TimeZone    *string
}
type MainStats struct {
	CurrentActivityName string
	TodayTracked        time.Duration
	TodaySessions       int
	StreakDays          int
}
