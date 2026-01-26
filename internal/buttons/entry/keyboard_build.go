package entry

import (
	"tracker-bot/pkg/buttonbuilder"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Reply button menus

func EntryReplyMenu() tgbotapi.ReplyKeyboardMarkup {
	return buttonbuilder.RK(
		buttonbuilder.RR(
			buttonbuilder.RB(EntryButtonProfile),
			buttonbuilder.RB(EntryButtonTrack),
		),
		buttonbuilder.RR(
			buttonbuilder.RB(EntryButtonLearning),
			buttonbuilder.RB(EntryButtonSubscription),
		),
	)
}
