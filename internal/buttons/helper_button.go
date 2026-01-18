// helpers for build telegram buttons and keyboardss
package buttonbuilder

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

// --- Inline helpers ----

// Inline buttons
func IB(text, data string) tgbotapi.InlineKeyboardButton {
	return tgbotapi.NewInlineKeyboardButtonData(text, data)
}

// inline button row
func IR(btns ...tgbotapi.InlineKeyboardButton) []tgbotapi.InlineKeyboardButton {
	return tgbotapi.NewInlineKeyboardRow(btns...)
}

// inline keyboard markup
func IK(rows ...[]tgbotapi.InlineKeyboardButton) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

// --- Replys helpers ----

// Reply buttons
func RB(text string) tgbotapi.KeyboardButton {
	return tgbotapi.NewKeyboardButton(text)
}

// reply button row
func RR(btns ...tgbotapi.KeyboardButton) []tgbotapi.KeyboardButton {
	return tgbotapi.NewKeyboardButtonRow(btns...)
}

// reply keyboard markup
func RK(rows ...[]tgbotapi.KeyboardButton) tgbotapi.ReplyKeyboardMarkup {
	kb := tgbotapi.NewReplyKeyboard(rows...)
	kb.ResizeKeyboard = true
	return kb
}
