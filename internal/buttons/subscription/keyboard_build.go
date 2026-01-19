package subscription

import (
	buttonbuilder "tracker-bot/internal/buttons"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Inline button menus

func SubscriptionEntryInlineMenu() tgbotapi.InlineKeyboardMarkup {
	return buttonbuilder.IK(
		buttonbuilder.IR(
			buttonbuilder.IB(SubscriptionButtonTariffPlans, SubscriptionCBTariffPlans),
			buttonbuilder.IB(SubscriptionButtonFreePlan, SubscriptionCBFreePlan),
		),
		buttonbuilder.IR(
			buttonbuilder.IB(SubscriptionButtonSupport, SubscriptionCBSupport),
			buttonbuilder.IB(SubscriptionButtonPaymentChange, SubscriptionCBPaymentChange),
		),
	)
}

// Reply button menus
