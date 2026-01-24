package dispatcher

import (
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
)

// pre-routing
func (d *Dispatcher) RunCallback(callback *tgbotapi.CallbackQuery) {
	if callback == nil || callback.Message == nil {
		log.Warn().Msg("CallbackQuery: nil callback tgbotapi.CallbackQuery")
		return

	} else {
		log.Debug().
			Str("user", fmt.Sprint(callback.From.ID)). //remake me
			Str("data", callback.Data).
			Msg("Callback context initialized")
	}
}

// Событие – нажатие инлайн-кнопки
func (d *Dispatcher) NewCallbackContext(cb *tgbotapi.CallbackQuery) *context.CallbackContext {

	ctxCB := &context.CallbackContext{
		Callback: cb,
		UserID:   int64(cb.From.ID), // Telegram ID
		Data:     cb.Data,
		Message:  cb.Message,
		UserName: strings.TrimSpace(cb.From.UserName),
	}
	if cb.Message != nil {
		ctxCB.ChatID = cb.Message.Chat.ID
	}

	reqCtx, cancel := stdctx.WithTimeout(stdctx.Background(), 3*time.Second)
	defer cancel()

	id, err := d.repo.EnsureIDByTelegram(reqCtx, ctxCB.UserID, ctxCB.UserName)
	if err != nil {
		log.Error().Err(err).Int64("tg_id", ctxCB.UserID).Str("data", cb.Data).Msg("ensure user failed (callback)")
		d.bot.Send(tgbotapi.NewCallback(cb.ID, "Ошибка, попробуйте ещё раз"))
		return ctxCB
	}

	ctxCB.DBUserID = id // users.id
	return ctxCB
}
