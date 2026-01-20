package subscription

import (
	"fmt"
	"tracker-bot/internal/models"
)

func SubscriptionMenuText(stats models.SubscriptionStats) string {
	return fmt.Sprintf(
		"%s\n\n%s *%s*\n%d *%d*\n%d *\n",
		SubscriptionUIMainTitle,
		SubscriptionUIMainTariffPlan, stats.ActivePlan,
		SubscriptionUIMainDaysEnd, stats.DaysEnd,
		SubscriptionUIMainMessage,
	)
}
