package router

import (
	"tracker-bot/internal/buttons/track"
	"tracker-bot/internal/service"
	tgclient "tracker-bot/internal/utils/tgcient"
	"tracker-bot/internal/utils/tgctx"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
)

//inline and replay routers

type Handler interface {
	Track()
}

type Module struct {
	bot      tgclient.TgBotAPI
	tracksvc *service.TrackerService
	entrysvc *service.EntryService
}

func New(bot tgclient.TgBotAPI, tracksvc *service.TrackerService, entrysvc *service.EntryService) *Module {
	return &Module{
		bot:      bot,
		tracksvc: tracksvc,
		entrysvc: entrysvc,
	}
}
func (m *Module) ShowTrackingMenu(ctx *tgctx.Message) {
	stats, err := m.tracksvc.GetMainStats(ctx.Ctx, ctx.DBUserID)
	if err != nil {
		log.Error().Err(err).Msg("GetMainStats failed")
		msg := tgbotapi.NewMessage(ctx.ChatID, "⚠️ Failed to load tracking data. Please try again.")
		_, _ = m.bot.Send(msg)
		return
	}

	text := track.TrackingMenuText(stats)

	msg := tgbotapi.NewMessage(ctx.ChatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = track.TrackEntryInlineMenu()

	if _, err := m.bot.Send(msg); err != nil {
		log.Error().Err(err).Msg("send tracking menu failed")
	}
}
