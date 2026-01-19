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

type ProfileStats struct {
	TgUserID    int64
	UserName    string
	PhoneNumber string
	Email       string
	Language    string
	TimeZone    string
}
type MainStats struct {
	CurrentActivityName string
	TodayTracked        time.Duration
	TodaySessions       int
	StreakDays          int
}

type LearningStats struct {
	Language     string
	TotalWords   int
	TodayWords   int
	LearnedWords int
	NextWordIn   string
}

type SubscriptionStats struct {
	ActivePlan string
	DaysEnd    int
}
