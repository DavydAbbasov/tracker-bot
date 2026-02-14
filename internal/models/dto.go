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

type TrackActivityItem struct {
	ID       int64
	Name     string
	Emoji    string
	Selected bool
}

type TimerDueUser struct {
	DBUserID    int64
	TgUserID    int64
	IntervalMin int
}

type ActivityDurationStat struct {
	ActivityID int64
	Name       string
	Emoji      string
	Duration   time.Duration
	Sessions   int
}

type ReportTodayStats struct {
	TotalTracked  time.Duration
	TotalSessions int
	TopActivities []ActivityDurationStat
}

type ReportPeriodStats struct {
	From          time.Time
	To            time.Time
	TotalTracked  time.Duration
	TotalSessions int
	Activities    []ActivityDurationStat
	Monthly       []MonthDurationStat
}

type MonthDurationStat struct {
	Month    time.Time
	Duration time.Duration
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
