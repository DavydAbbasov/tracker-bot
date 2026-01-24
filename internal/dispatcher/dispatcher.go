package dispatcher

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"

	h "tracker-bot/internal/buttons/handlers"
	router "tracker-bot/internal/handlers"
	tgclient "tracker-bot/internal/utils/tgcient"
	"tracker-bot/internal/utils/tgctx"
)

type Dispatcher struct {
	bot          tgclient.TgBotAPI
	track        *router.Module
	subscription *router.Module
	entry        *router.Module
	profile      *router.Module
	learning     *router.Module

	reply *h.ReplyModule
}

func New(
	bot tgclient.TgBotAPI,
	track *router.Module,
	subscription *router.Module,
	entry *router.Module,
	profile *router.Module,
	learning *router.Module,
) *Dispatcher {
	if bot == nil {
		log.Fatal().Msg("Dispatcher: nil bot interfaces.BotAPI")
	}

	d := &Dispatcher{
		bot:          bot,
		track:        track,
		subscription: subscription,
		entry:        entry,
		profile:      profile,
		learning:     learning,
	}

	d.reply = h.New(bot, track, subscription, entry, profile, learning)

	return d
}

func (d *Dispatcher) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := d.bot.GetUpdatesChan(u)

	for update := range updates {
		switch {
		case update.CallbackQuery != nil:
			d.RunCallback(update.CallbackQuery)
		case update.Message != nil:
			d.handleMessage(update.Message)
		}
	}
}

func (d *Dispatcher) handleMessage(msg *tgbotapi.Message) {
	ctx := d.newMessageContext(msg)

	if d.handleUserState(ctx) {
		return
	}

	if d.reply.HandleReplyButtons(ctx) {
		return
	}

	if msg.IsCommand() {
		d.handleCommand(ctx)
		return
	}

	d.handleText(ctx)
}

func (d *Dispatcher) newMessageContext(msg *tgbotapi.Message) *tgctx.MsgContext {

	ctxMsg := &tgctx.MsgContext{
		ChatID: msg.Chat.ID,
		UserID: int64(msg.From.ID), // телеграм-ID (почему берем под int64)
		Text:   msg.Text,
	}

	return ctxMsg
}
func (d *Dispatcher) handleUserState(ctx *tgctx.MsgContext) bool {
	return false

}

// примеры, чтобы было понятно, что делать дальше:
func (d *Dispatcher) handleCommand(ctx *tgctx.MsgContext) {
	// ctx.Message.Command() если ты сохраняешь Message в ctx
	// или ctx.Text если там "/start"
}

func (d *Dispatcher) handleText(ctx *tgctx.MsgContext) {
	// fallback логика
}
