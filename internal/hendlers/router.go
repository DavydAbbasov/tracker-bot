package router

import (
	"tracker-bot/internal/buttons/learning"
	"tracker-bot/internal/buttons/profile"
	"tracker-bot/internal/buttons/subscription"
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
	bot             tgclient.TgBotAPI
	profilesvc      *service.ProfileService
	tracksvc        *service.TrackerService
	learningsvc     *service.LearningService
	subscriptionsvc *service.SubscriptionService
	entrysvc        *service.EntryService
}

func New(bot tgclient.TgBotAPI, profilesvc *service.ProfileService, tracksvc *service.TrackerService, learningsvc *service.LearningService, subscriptionsvc *service.SubscriptionService, entrysvc *service.EntryService) *Module {
	return &Module{
		bot:             bot,
		profilesvc:      profilesvc,
		tracksvc:        tracksvc,
		learningsvc:     learningsvc,
		subscriptionsvc: subscriptionsvc,
		entrysvc:        entrysvc,
	}
}

func (m *Module) ShowProfileMenu(ctx *tgctx.MsgContext) {
	stats, err := m.profilesvc.GetProfileStats(ctx.Ctx, ctx.DBUserID)
	if err != nil {
		log.Error().Err(err).Msg("GetProfile failed")
		msg := tgbotapi.NewMessage(ctx.ChatID, "⚠️ Failed to load profile data. Please try again.")
		_, _ = m.bot.Send(msg)
		return
	}

	text := profile.ProfileMenuText(stats)

	msg := tgbotapi.NewMessage(ctx.ChatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = profile.ProfileEntryInlineMenu()

	if _, err := m.bot.Send(msg); err != nil {
		log.Error().Err(err).Msg("send profile menu failed")
	}
}

func (m *Module) ShowTrackingMenu(ctx *tgctx.MsgContext) {
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

func (m *Module) ShowLearningMenu(ctx *tgctx.MsgContext) {
	stats, err := m.learningsvc.GetLearningStats(ctx.Ctx, ctx.DBUserID)
	if err != nil {
		log.Error().Err(err).Msg("GetLearningStats failed")
		msg := tgbotapi.NewMessage(ctx.ChatID, "⚠️ Failed to load learning data. Please try again.")
		_, _ = m.bot.Send(msg)
		return
	}

	text := learning.LearningMenuText(stats)

	msg := tgbotapi.NewMessage(ctx.ChatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = learning.LearningEntryInlineMenu()

	if _, err := m.bot.Send(msg); err != nil {
		log.Error().Err(err).Msg("send learning menu failed")
	}
}

func (m *Module) ShowSubscriptionMenu(ctx *tgctx.MsgContext) {
	stats, err := m.subscriptionsvc.GetSubscriptionStats(ctx.Ctx, ctx.DBUserID)
	if err != nil {
		log.Error().Err(err).Msg("GetSubscriptionStats failed")
		msg := tgbotapi.NewMessage(ctx.ChatID, "⚠️ Failed to load subscription data. Please try again.")
		_, _ = m.bot.Send(msg)
		return
	}

	text := subscription.SubscriptionInlineMenuText(stats)

	msg := tgbotapi.NewMessage(ctx.ChatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = subscription.SubscriptionEntryInlineMenu()

	if _, err := m.bot.Send(msg); err != nil {
		log.Error().Err(err).Msg("send subscription menu failed")
	}
}
