package track

import (
	"fmt"
	"strings"
	"tracker-bot/internal/models"
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
		buttonbuilder.RR(buttonbuilder.RB(TrackButtonActivityActivate)),
		buttonbuilder.RR(buttonbuilder.RB(TrackButtonActivityDelete), buttonbuilder.RB(TrackButtonViewArchive)),
		buttonbuilder.RR(buttonbuilder.RB(TrackButtonBackHome)),
	)
}

func TrackArchiveReplyMenu() tgbotapi.ReplyKeyboardMarkup {
	return buttonbuilder.RK(
		buttonbuilder.RR(buttonbuilder.RB(TrackButtonSelectActivity), buttonbuilder.RB(TrackButtonViewArchive)),
		buttonbuilder.RR(buttonbuilder.RB(TrackButtonBackHome)),
	)
}

func TrackTimerReplyMenu() tgbotapi.ReplyKeyboardMarkup {
	return buttonbuilder.RK(
		buttonbuilder.RR(buttonbuilder.RB(TrackButtonTimer15), buttonbuilder.RB(TrackButtonTimer30)),
		buttonbuilder.RR(buttonbuilder.RB(TrackButtonBackHome)),
	)
}

func TrackActivitiesInlineMenu(items []models.TrackActivityItem) tgbotapi.InlineKeyboardMarkup {
	rows := make([][]tgbotapi.InlineKeyboardButton, 0, len(items)+1)
	for _, item := range items {
		if strings.TrimSpace(item.Name) == "" {
			continue
		}

		check := "⚪"
		if item.Selected {
			check = "🟢"
		}

		title := check + " " + item.Name
		if item.Emoji != "" {
			title = check + " " + item.Emoji + " " + item.Name
		}

		callbackData := fmt.Sprintf("act_toggle_:%d", item.ID)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(title, callbackData),
		))
	}

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("🛒 Archive selected", TrackCBArchiveSelected),
	))
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("↩️ Back", "back_to_main"),
	))

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func TrackPromptInlineMenu(items []models.TrackActivityItem, intervalMin int) tgbotapi.InlineKeyboardMarkup {
	rows := make([][]tgbotapi.InlineKeyboardButton, 0, len(items)+1)
	for _, item := range items {
		if strings.TrimSpace(item.Name) == "" {
			continue
		}
		title := item.Name
		if item.Emoji != "" {
			title = item.Emoji + " " + item.Name
		}
		callbackData := fmt.Sprintf("%s%d:%d", TrackCBPromptActivity, item.ID, intervalMin)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(title, callbackData),
		))
	}
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("⏹ Stop Timer", TrackCBPromptStopTimer),
	))
	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func TrackArchiveInlineMenu(items []models.TrackActivityItem) tgbotapi.InlineKeyboardMarkup {
	rows := make([][]tgbotapi.InlineKeyboardButton, 0, len(items)*2+1)
	for _, item := range items {
		if strings.TrimSpace(item.Name) == "" {
			continue
		}
		title := item.Name
		if item.Emoji != "" {
			title = item.Emoji + " " + item.Name
		}
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📦 "+title, "noop"),
		))
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("♻ Restore", fmt.Sprintf("%s%d", TrackCBArchiveRestore, item.ID)),
			tgbotapi.NewInlineKeyboardButtonData("🗑 Delete forever", fmt.Sprintf("%s%d", TrackCBArchiveDelete, item.ID)),
		))
	}
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("📂 Active activities", TrackCBArchiveToActive),
	))
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("↩️ Back", "back_to_main"),
	))
	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func TrackCreateSuccessInlineMenu() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📂 Open Activities", TrackCBOpenActivities),
			tgbotapi.NewInlineKeyboardButtonData("➕ Create Another", TrackCBCreateAnother),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("↩️ Back", "back_to_main"),
		),
	)
}

func TrackArchiveSuccessInlineMenu() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🗄 Open Archive", TrackCBOpenArchive),
			tgbotapi.NewInlineKeyboardButtonData("📂 Open Activities", TrackCBOpenActivities),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("↩️ Back", "back_to_main"),
		),
	)
}
