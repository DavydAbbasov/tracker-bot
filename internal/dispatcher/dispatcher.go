package dispatcher

import (
	"context"
	"tracker-bot/internal/models"
	"tracker-bot/internal/service"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"

	h "tracker-bot/internal/buttons/handlers"
	router "tracker-bot/internal/handlers"
	"tracker-bot/internal/utils/tgctx"
)

type Dispatcher struct {
	bot          *tgbotapi.BotAPI
	appCtx       context.Context
	entrysvc     service.EntryService
	track        *router.Module
	subscription *router.Module
	entry        *router.Module
	profile      *router.Module
	learning     *router.Module

	reply *h.ReplyModule
}

func New(
	bot *tgbotapi.BotAPI,
	appCtx context.Context,
	entrysvc service.EntryService,
	track *router.Module,
	subscription *router.Module,
	entry *router.Module,
	profile *router.Module,
	learning *router.Module,
) *Dispatcher {
	if bot == nil {
		log.Fatal().Msg("Dispatcher: nil bot interfaces.BotAPI")
	}

	if appCtx == nil {
		appCtx = context.Background()
	}

	d := &Dispatcher{
		bot:          bot,
		appCtx:       appCtx,
		entrysvc:     entrysvc,
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

func (d *Dispatcher) ensureUser(ctx *tgctx.MsgContext, chatID int64, from *tgbotapi.User) bool {
	if from == nil {
		return false
	}

	in := &models.UserInput{
		TgUserID: int64(from.ID),
		UserName: &from.UserName,
	}

	dbID, err := d.entrysvc.EnsureUser(ctx.Ctx, in)
	if err != nil {
		log.Error().Err(err).Msg("ensure user failed")
		out := tgbotapi.NewMessage(chatID, "⚠️ Ошибка. Попробуй ещё раз.")
		_, _ = d.bot.Send(out)
		return false
	}
	ctx.DBUserID = dbID
	return true
}

func (d *Dispatcher) newMessageContext(msg *tgbotapi.Message) *tgctx.MsgContext {
	ctx := &tgctx.MsgContext{
		Ctx:    d.appCtx,
		ChatID: msg.Chat.ID,
		Text:   msg.Text,
	}

	if msg.From != nil {
		ctx.UserID = int64(msg.From.ID)
	}

	return ctx
}

func (d *Dispatcher) handleMessage(msg *tgbotapi.Message) {
	if msg == nil || msg.From == nil {
		return
	}

	mctx := d.newMessageContext(msg)

	if !d.ensureUser(mctx, msg.Chat.ID, msg.From) {
		return
	}

	// 1) команды СНАЧАЛА (чтобы /start не шёл в reply)
	if msg.IsCommand() {
		d.handleCommand(msg, mctx)
		return
	}

	// 2) FSM
	if d.handleUserState(mctx) {
		return
	}

	// 3) reply-кнопки
	if d.reply != nil && d.reply.HandleReplyButtons(mctx) {
		return
	}

	// 4) обычный текст
	d.handleText(mctx)
}

func (d *Dispatcher) handleCallback(q *tgbotapi.CallbackQuery) {
	if q == nil || q.Message == nil || q.From == nil {
		return
	}

	ack := tgbotapi.NewCallback(q.ID, "")
	if _, err := d.bot.Request(ack); err != nil {
		log.Error().Err(err).Msg("callback ack failed")
	}

	mctx := &tgctx.MsgContext{
		Ctx:    d.appCtx,
		ChatID: q.Message.Chat.ID,
		Text:   q.Data,
		UserID: int64(q.From.ID),
	}

	if !d.ensureUser(mctx, q.Message.Chat.ID, q.From) {
		return
	}

	if d.reply != nil && d.reply.HandleReplyButtons(mctx) {
		return
	}
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
