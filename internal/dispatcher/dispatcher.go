package dispatcher

import (
	"context"
	"strings"
	trackbtn "tracker-bot/internal/buttons/track"
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

	waitingActivityName map[int64]bool
	userScreen          map[int64]string
}

const (
	screenUnknown        = ""
	screenHome           = "home"
	screenTrackMain      = "track_main"
	screenTrackManage    = "track_manage"
	screenTrackTimer     = "track_timer"
	screenTrackArchive   = "track_archive"
	screenCreateActivity = "create_activity"
)

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
		bot:                 bot,
		appCtx:              appCtx,
		entrysvc:            entrysvc,
		track:               track,
		subscription:        subscription,
		entry:               entry,
		profile:             profile,
		learning:            learning,
		waitingActivityName: make(map[int64]bool),
		userScreen:          make(map[int64]string),
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
	if ctxText := mctx.Text; ctxText == "📈Track" {
		d.userScreen[mctx.UserID] = screenTrackMain
	}
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
		Ctx:       d.appCtx,
		ChatID:    q.Message.Chat.ID,
		Text:      q.Data,
		UserID:    int64(q.From.ID),
		MessageID: q.Message.MessageID,
	}

	if !d.ensureUser(mctx, q.Message.Chat.ID, q.From) {
		return
	}

	if strings.HasPrefix(q.Data, "track:") || strings.HasPrefix(q.Data, "act_toggle_:") {
		d.handleTrackCallback(mctx, q.Data)
		return
	}

	if d.reply != nil && d.reply.HandleReplyButtons(mctx) {
		return
	}
}

func (d *Dispatcher) handleUserState(ctx *tgctx.MsgContext) bool {
	if d.waitingActivityName[ctx.UserID] {
		if d.isTrackButtonText(ctx.Text) {
			_, _ = d.bot.Send(tgbotapi.NewMessage(ctx.ChatID, "Use buttons from menu. Enter activity name as plain text."))
			return true
		}
		done := d.track.ProcessCreateActivity(ctx)
		if done {
			delete(d.waitingActivityName, ctx.UserID)
			d.userScreen[ctx.UserID] = screenTrackMain
		}
		return true
	}

	return false
}

func (d *Dispatcher) handleCommand(msg *tgbotapi.Message, ctx *tgctx.MsgContext) {
	cmd := msg.Command()

	switch cmd {
	case "start":
		d.userScreen[ctx.UserID] = screenHome
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
	switch ctx.Text {
	case trackbtn.TrackButtonActivityDelete:
		if d.userScreen[ctx.UserID] != screenTrackManage {
			d.replyUseButtons(ctx.ChatID)
			return
		}
		d.track.DeleteSelectedActivities(ctx)
		return
	case trackbtn.TrackButtonActivityActivate:
		if d.userScreen[ctx.UserID] != screenTrackManage && d.userScreen[ctx.UserID] != screenTrackMain {
			d.replyUseButtons(ctx.ChatID)
			return
		}
		d.userScreen[ctx.UserID] = screenTrackTimer
		d.track.ShowTrackTimerMenu(ctx)
		return
	case trackbtn.TrackButtonActivityArchive:
		d.userScreen[ctx.UserID] = screenTrackArchive
		d.track.ShowArchiveMenu(ctx)
		return
	case trackbtn.TrackButtonViewArchive:
		d.userScreen[ctx.UserID] = screenTrackArchive
		d.track.ShowArchiveMenu(ctx)
		return
	case trackbtn.TrackButtonSelectActivity:
		d.userScreen[ctx.UserID] = screenTrackManage
		d.track.ShowTrackActivitySelectionMenu(ctx)
		return
	case trackbtn.TrackButtonTimer15:
		if d.userScreen[ctx.UserID] != screenTrackTimer {
			d.replyUseButtons(ctx.ChatID)
			return
		}
		d.track.ActivateTrackTimer(ctx, 15)
		d.userScreen[ctx.UserID] = screenHome
		return
	case trackbtn.TrackButtonTimer30:
		if d.userScreen[ctx.UserID] != screenTrackTimer {
			d.replyUseButtons(ctx.ChatID)
			return
		}
		d.track.ActivateTrackTimer(ctx, 30)
		d.userScreen[ctx.UserID] = screenHome
		return
	case trackbtn.TrackButtonBackHome:
		d.userScreen[ctx.UserID] = screenHome
		d.entry.ShowEntryMenu(ctx)
		return
	}

	out := tgbotapi.NewMessage(ctx.ChatID, "Я тебя понял, но не знаю что с этим сделать. Напиши /help")
	if _, err := d.bot.Send(out); err != nil {
		log.Error().Err(err).Msg("send fallback failed")
	}
}

func (d *Dispatcher) handleTrackCallback(ctx *tgctx.MsgContext, data string) {
	switch {
	case data == "noop":
		return
	case data == "back_to_main":
		d.userScreen[ctx.UserID] = screenTrackMain
		d.track.ShowTrackingMenu(ctx)
	case data == trackbtn.TrackCBActivitySelect:
		d.userScreen[ctx.UserID] = screenTrackManage
		d.track.ShowTrackActivitySelectionMenu(ctx)
	case data == trackbtn.TrackCBActivityCreate:
		d.waitingActivityName[ctx.UserID] = true
		d.userScreen[ctx.UserID] = screenCreateActivity
		d.track.PromptCreateActivity(ctx)
	case data == trackbtn.TrackCBArchiveOpen:
		d.userScreen[ctx.UserID] = screenTrackArchive
		d.track.ShowArchiveMenu(ctx)
	case data == trackbtn.TrackCBOpenArchive:
		d.userScreen[ctx.UserID] = screenTrackArchive
		d.track.ShowArchiveMenuInPlace(ctx)
	case data == trackbtn.TrackCBOpenActivities:
		d.userScreen[ctx.UserID] = screenTrackManage
		d.track.ShowTrackActivitySelectionMenuInPlace(ctx)
	case data == trackbtn.TrackCBCreateAnother:
		d.waitingActivityName[ctx.UserID] = true
		d.userScreen[ctx.UserID] = screenCreateActivity
		d.track.PromptCreateActivity(ctx)
	case data == trackbtn.TrackCBArchiveSelected:
		d.userScreen[ctx.UserID] = screenTrackArchive
		d.track.ArchiveSelectedActivitiesInPlace(ctx)
	case data == trackbtn.TrackCBArchiveToActive:
		d.userScreen[ctx.UserID] = screenTrackManage
		d.track.ShowTrackActivitySelectionMenuInPlace(ctx)
	case data == trackbtn.TrackCBPromptStopTimer:
		d.track.StopTrackTimer(ctx)
	case strings.HasPrefix(data, trackbtn.TrackCBPromptActivity):
		d.track.RecordPromptAnswer(ctx)
	case strings.HasPrefix(data, trackbtn.TrackCBArchiveRestore):
		d.track.RestoreArchivedActivity(ctx)
	case strings.HasPrefix(data, trackbtn.TrackCBArchiveDelete):
		d.track.DeleteArchivedForever(ctx)
	case strings.HasPrefix(data, "act_toggle_:"):
		d.track.HandleTrackToggleCallback(ctx)
	}
}

func (d *Dispatcher) replyUseButtons(chatID int64) {
	_, _ = d.bot.Send(tgbotapi.NewMessage(chatID, "Use buttons from menu."))
}

func (d *Dispatcher) isTrackButtonText(text string) bool {
	switch text {
	case trackbtn.TrackButtonActivityActivate,
		trackbtn.TrackButtonActivityArchive,
		trackbtn.TrackButtonActivityDelete,
		trackbtn.TrackButtonTimer15,
		trackbtn.TrackButtonTimer30,
		trackbtn.TrackButtonBackHome,
		trackbtn.TrackButtonViewArchive:
		return true
	default:
		return false
	}
}
