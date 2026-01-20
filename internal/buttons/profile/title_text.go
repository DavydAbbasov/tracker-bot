package profile

import (
	"fmt"
	"tracker-bot/internal/models"
)

func ProfileMenuText(stats models.ProfileStats) string {
	return fmt.Sprintf(
		"%s\n\n%d *%s*\n%s *%s*\n%s *\n",
		ProfileUIMainTitle,
		ProfileUIMainID, stats.TgUserID,
		ProfileUIMainName, stats.UserName,
		ProfileUIMainLanguage, stats.Language,
		ProfileUIMainTimeZone, stats.TimeZone,
		ProfileUIMainEmail, stats.Email,
	)
}
