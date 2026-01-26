package track

import (
	"tracker-bot/pkg/buttonbuilder"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Inline botton menus

func TrackEntryInlineMenu() tgbotapi.InlineKeyboardMarkup {
	return buttonbuilder.IK(
		buttonbuilder.IR(
			buttonbuilder.IB(TrackButtonSelectActivity, TrackCBActivitySelect),
			buttonbuilder.IB(TrackButtonCreateActivity, TrackCBActivityCreate),
		),
		buttonbuilder.IR(
			buttonbuilder.IB(TrackButtonViewReports, TrackCBReportSummary),
			buttonbuilder.IB(TrackButtonViewArchive, TrackCBArchiveOpen),
		),
	)
}

// Reply button menus

func TrackActivityListReplyMenu() tgbotapi.ReplyKeyboardMarkup {
	return buttonbuilder.RK(
		buttonbuilder.RR(buttonbuilder.RB(TrackButtonToday), buttonbuilder.RB(TrackButtonBack)),
	)
}

func TrackActivityReportReplyMenu() tgbotapi.ReplyKeyboardMarkup {
	return buttonbuilder.RK(
		buttonbuilder.RR(buttonbuilder.RB(TrackButtonReportPeriod), buttonbuilder.RB(TrackButtonReportWeek)),
		buttonbuilder.RR(buttonbuilder.RB(TrackButtonReportExport), buttonbuilder.RB(TrackButtonToday)),
		buttonbuilder.RR(buttonbuilder.RB(TrackButtonReportDelete), buttonbuilder.RB(TrackButtonBack)),
	)
}

func TrackActivityManageReplyMenu() tgbotapi.ReplyKeyboardMarkup {
	return buttonbuilder.RK(
		buttonbuilder.RR(buttonbuilder.RB(TrackButtonActivityActivate), buttonbuilder.RB(TrackButtonActivityArchive)),
		buttonbuilder.RR(buttonbuilder.RB(TrackButtonActivityDelete), buttonbuilder.RB(TrackButtonBackHome)),
	)
}

func TrackTimerReplyMenu() tgbotapi.ReplyKeyboardMarkup {
	return buttonbuilder.RK(
		buttonbuilder.RR(buttonbuilder.RB(TrackButtonTimer15), buttonbuilder.RB(TrackButtonTimer60)),
		buttonbuilder.RR(buttonbuilder.RB(TrackButtonTimerCreate), buttonbuilder.RB(TrackButtonBackHome)),
	)
}
