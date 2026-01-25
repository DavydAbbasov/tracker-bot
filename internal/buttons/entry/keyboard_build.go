package entry

import (
	buttonbuilder "tracker-bot/internal/buttons"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Inline button menus

func EntryInlineMenu() tgbotapi.InlineKeyboardMarkup {
	return buttonbuilder.IK(
		buttonbuilder.IR(
			buttonbuilder.IB(EntryButtonProfile, EntryCBProfile),
			buttonbuilder.IB(EntryButtonTrack, EntryCBTrack),
		),
		buttonbuilder.IR(
			buttonbuilder.IB(EntryButtonLearning, EntryCBLearning),
			buttonbuilder.IB(EntryButtonSubscription, EntryCBSubscription),
		),
	)
}
