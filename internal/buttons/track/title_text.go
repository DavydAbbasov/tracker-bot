package track

import (
	"fmt"
	"time"
	"tracker-bot/internal/models"
)

// Main screen
type TrackMainStats struct {
	CurrentActivityName string
	TodayTrackedTime    string
	TodayActivityCount  string
	CurrentStreakDays   string
}

// Activity report
type TrackActivityReportStats struct {
	ActivityStartDate    string
	ConsecutiveDaysCount string
	TodayAccumulatedTime string
	AverageDailyTime     string
	ReportDate           string
}

func TrackingMenuText(stats models.MainStats) string {
	return fmt.Sprintf(
		"%s\n\n%s *%s*\n%s *%s*\n%s *%d*\n%s *%d*\n",
		TrackUIMainTitle,
		TrackUIMainLabelCurrentActivity, safeText(stats.CurrentActivityName),
		TrackUIMainLabelTodayTime, formatDuration(stats.TodayTracked),
		TrackUIMainLabelStreak, stats.StreakDays,
		TrackUIMainLabelTodayCount, stats.TodaySessions,
	)
}

// formatDuration formats duration into human-readable string like "4h 30m".
func formatDuration(d time.Duration) string {
	if d < 0 {
		d = 0
	}
	h := int(d.Hours())
	m := int(d.Minutes()) % 60

	switch {
	case h > 0 && m > 0:
		return fmt.Sprintf("%dh %dm", h, m)
	case h > 0:
		return fmt.Sprintf("%dh", h)
	default:
		return fmt.Sprintf("%dm", m)
	}
}

// safeText returns fallback when string is empty.
func safeText(s string) string {
	if s == "" {
		return "—"
	}
	return s
}
