package dispatcher

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"

	h "tracker-bot/internal/buttons/handlers"
	router "tracker-bot/internal/handlers"
	"tracker-bot/internal/utils/tgctx"
)

type Dispatcher struct {
	bot          *tgbotapi.BotAPI
	track        *router.Module
	subscription *router.Module
	entry        *router.Module
	profile      *router.Module
	learning     *router.Module

	reply *h.ReplyModule
}

func New(
	bot *tgbotapi.BotAPI,
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
		case update.Message != nil:
			d.handleMessage(update.Message)

		case update.CallbackQuery != nil:
			d.handleCallback(update.CallbackQuery)
		}
	}
}

func (d *Dispatcher) handleMessage(msg *tgbotapi.Message) {
	if msg == nil {
		return
	}

	ctx := d.newMessageContext(msg)

	// 1) обработка состояний пользователя (если появится FSM)
	if d.handleUserState(ctx) {
		return
	}

	// 2) reply-кнопки
	if d.reply != nil && d.reply.HandleReplyButtons(ctx) {
		return
	}

	// 3) команды
	if msg.IsCommand() {
		d.handleCommand(msg, ctx)
		return
	}

	// 4) обычный текст
	d.handleText(ctx)
}

func (d *Dispatcher) handleCallback(q *tgbotapi.CallbackQuery) {
	if q == nil || q.Message == nil {
		return
	}

	ack := tgbotapi.NewCallback(q.ID, "")
	if _, err := d.bot.Request(ack); err != nil {
		log.Error().Err(err).Msg("callback ack failed")
	}

	ctx := &tgctx.MsgContext{
		ChatID: q.Message.Chat.ID,
		Text:   q.Data,
	}

	if q.From != nil {
		ctx.UserID = int64(q.From.ID)
	}

	if d.reply != nil && d.reply.HandleReplyButtons(ctx) {
		return
	}
}

func (d *Dispatcher) newMessageContext(msg *tgbotapi.Message) *tgctx.MsgContext {
	ctx := &tgctx.MsgContext{
		ChatID: msg.Chat.ID,
		Text:   msg.Text,
	}

	if msg.From != nil {
		ctx.UserID = int64(msg.From.ID)
	}

	return ctx
}

func (d *Dispatcher) handleUserState(ctx *tgctx.MsgContext) bool {
	return false
}

func (d *Dispatcher) handleCommand(msg *tgbotapi.Message, ctx *tgctx.MsgContext) {
	cmd := msg.Command()

	switch cmd {
	case "start":
		d.entry.ShowEntryMenu(ctx)
		return

	case "help":
		out := tgbotapi.NewMessage(ctx.ChatID, "Доступные команды: /start, /help")
		if _, err := d.bot.Send(out); err != nil {
			log.Error().Err(err).Msg("send help failed")
		}
		return

	default:
		out := tgbotapi.NewMessage(ctx.ChatID, "Неизвестная команда.")
		if _, err := d.bot.Send(out); err != nil {
			log.Error().Err(err).Msg("send unknown command failed")
		}
		return
	}
}

func (d *Dispatcher) handleText(ctx *tgctx.MsgContext) {
	out := tgbotapi.NewMessage(ctx.ChatID, "Я тебя понял, но не знаю что с этим сделать. Напиши /help")
	if _, err := d.bot.Send(out); err != nil {
		log.Error().Err(err).Msg("send fallback failed")
	}
}
