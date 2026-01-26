package profile

import (
	"fmt"
	"tracker-bot/internal/models"
	"tracker-bot/pkg/textbuilder"
)

func ProfileMenuText(stats *models.ProfileStats) string {
	return fmt.Sprintf(
		"%s\n\n"+
			"%s %d\n"+
			"%s %s\n"+
			"%s %s\n"+
			"%s %s\n"+
			"%s %s",
		ProfileUIMainTitle,
		ProfileUIMainID, stats.TgUserID,
		ProfileUIMainName, textbuilder.StrOrDashMD(stats.UserName),
		ProfileUIMainLanguage, textbuilder.StrOrDashMD(stats.Language),
		ProfileUIMainTimeZone, textbuilder.StrOrDashMD(stats.TimeZone),
		ProfileUIMainEmail, textbuilder.StrOrDashMD(stats.Email),
	)
}
