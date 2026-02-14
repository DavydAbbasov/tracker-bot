package handlers

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"
	"tracker-bot/internal/buttons/entry"
	"tracker-bot/internal/buttons/learning"
	"tracker-bot/internal/buttons/profile"
	"tracker-bot/internal/buttons/subscription"
	"tracker-bot/internal/buttons/track"
	"tracker-bot/internal/models"
	"tracker-bot/internal/service"
	"tracker-bot/internal/utils/tgctx"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
)

// Module routes UI actions to services and renders bot responses.

type Handler interface {
	Track()
}

type Module struct {
	bot             *tgbotapi.BotAPI
	profilesvc      service.ProfileService
	tracksvc        service.TrackerService
	timersvc        service.TimerService
	learningsvc     service.LearningService
	subscriptionsvc service.SubscriptionService
	entrysvc        service.EntryService
	testTimerMin    int
}

func New(bot *tgbotapi.BotAPI, entrysvc service.EntryService, profilesvc service.ProfileService, tracksvc service.TrackerService, timersvc service.TimerService, learningsvc service.LearningService, subscriptionsvc service.SubscriptionService, testTimerMin int) *Module {
	return &Module{
		bot:             bot,
		profilesvc:      profilesvc,
		tracksvc:        tracksvc,
		timersvc:        timersvc,
		learningsvc:     learningsvc,
		subscriptionsvc: subscriptionsvc,
		entrysvc:        entrysvc,
		testTimerMin:    testTimerMin,
	}
}

func (m *Module) ShowEntryMenu(ctx *tgctx.MsgContext) {
	text := entry.EntryMenuText()

	msg := tgbotapi.NewMessage(ctx.ChatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = entry.EntryReplyMenu()

	if _, err := m.bot.Send(msg); err != nil {
		log.Error().Err(err).Msg("send entry menu failed")
	}
}

func (m *Module) ShowProfileMenu(ctx *tgctx.MsgContext) {
	stats, err := m.profilesvc.GetProfileStats(ctx.Ctx, ctx.UserID)
	if err != nil {
		log.Error().Err(err).Msg("GetProfile failed")
		msg := tgbotapi.NewMessage(ctx.ChatID, "⚠️ Failed to load profile data. Please try again.")
		_, _ = m.bot.Send(msg)
		return
	}

	text := profile.ProfileMenuText(stats)

	msg := tgbotapi.NewMessage(ctx.ChatID, text)
	msg.ReplyMarkup = profile.ProfileEntryInlineMenu()

	if _, err := m.bot.Send(msg); err != nil {
		log.Error().Err(err).Msg("send profile menu failed")
	}
}

func (m *Module) ShowTrackingMenu(ctx *tgctx.MsgContext) {
	stats, err := m.tracksvc.GetMainStats(ctx.Ctx, ctx.UserID)
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

func (m *Module) ShowReportsHub(ctx *tgctx.MsgContext, inPlace bool) {
	text := "📈 Reports\n\nChoose a report type:"
	msgReply := tgbotapi.NewMessage(ctx.ChatID, "📈")
	msgReply.ReplyMarkup = track.TrackReportsReplyMenu()
	_, _ = m.bot.Send(msgReply)

	if inPlace && ctx.MessageID > 0 {
		msg := tgbotapi.NewEditMessageText(ctx.ChatID, ctx.MessageID, text)
		_, _ = m.bot.Send(msg)
		return
	}
	_, _ = m.bot.Send(tgbotapi.NewMessage(ctx.ChatID, text))
}

func (m *Module) ShowTodayChart(ctx *tgctx.MsgContext) {
	stats, err := m.tracksvc.GetTodayReport(ctx.Ctx, ctx.DBUserID)
	if err != nil {
		log.Error().Err(err).Msg("today chart failed")
		_, _ = m.bot.Send(tgbotapi.NewMessage(ctx.ChatID, "⚠️ Failed to load chart data."))
		return
	}
	if len(stats.TopActivities) == 0 {
		_, _ = m.bot.Send(tgbotapi.NewMessage(ctx.ChatID, "📉 No data for chart yet."))
		return
	}

	maxDur := time.Duration(1)
	for _, a := range stats.TopActivities {
		if a.Duration > maxDur {
			maxDur = a.Duration
		}
	}

	var b strings.Builder
	b.WriteString("📉 Today Chart\n\n")
	total := stats.TotalTracked
	for _, a := range stats.TopActivities {
		name := a.Name
		if a.Emoji != "" {
			name = a.Emoji + " " + a.Name
		}
		barLen := int((float64(a.Duration) / float64(maxDur)) * 12.0)
		if barLen < 1 {
			barLen = 1
		}
		if barLen > 12 {
			barLen = 12
		}
		percent := percentOf(a.Duration, total)
		b.WriteString(fmt.Sprintf("%s\n%s %s (%s)\n\n", name, strings.Repeat("█", barLen), formatReportDuration(a.Duration), percent))
	}

	msg := tgbotapi.NewMessage(ctx.ChatID, b.String())
	msg.ReplyMarkup = track.TrackReportTodayInlineMenu()
	_, _ = m.bot.Send(msg)
}

func (m *Module) ShowPeriodMenu(ctx *tgctx.MsgContext, selected map[int64]bool, month, from, to time.Time) {
	items, err := m.tracksvc.ListActivities(ctx.Ctx, ctx.DBUserID)
	if err != nil {
		_, _ = m.bot.Send(tgbotapi.NewMessage(ctx.ChatID, "⚠️ Failed to load activities for period."))
		return
	}
	if month.IsZero() {
		month = time.Now().UTC()
	}
	rangeLabel := formatDateOrDash(from) + ".." + formatDateOrDash(to)
	text := fmt.Sprintf("📅 Period Report\nSelected: %d activities\nRange: %s", len(selected), rangeLabel)
	if ctx.MessageID > 0 {
		edit := tgbotapi.NewEditMessageTextAndMarkup(
			ctx.ChatID,
			ctx.MessageID,
			text,
			track.TrackReportPeriodInlineMenu(items, selected, rangeLabel),
		)
		_, _ = m.bot.Send(edit)
		return
	}
	msg := tgbotapi.NewMessage(ctx.ChatID, text)
	msg.ReplyMarkup = track.TrackReportPeriodInlineMenu(items, selected, rangeLabel)
	_, _ = m.bot.Send(msg)
}

func (m *Module) ShowPeriodTextReport(ctx *tgctx.MsgContext, from, to time.Time, activityIDs []int64, selectedOnly bool) {
	stats, err := m.tracksvc.GetPeriodReport(ctx.Ctx, ctx.DBUserID, from, to.Add(24*time.Hour), activityIDs)
	if err != nil {
		_, _ = m.bot.Send(tgbotapi.NewMessage(ctx.ChatID, "⚠️ Failed to build period report."))
		return
	}
	var b strings.Builder
	b.WriteString("📄 Period Report\n\n")
	b.WriteString(fmt.Sprintf("Range: %s..%s\n", from.Format("2006-01-02"), to.Format("2006-01-02")))
	if selectedOnly {
		b.WriteString("Scope: selected activities\n")
	} else {
		b.WriteString("Scope: all selected in menu\n")
	}
	b.WriteString(fmt.Sprintf("Total: %s\nSessions: %d\n\n", formatReportDuration(stats.TotalTracked), stats.TotalSessions))
	total := stats.TotalTracked
	if len(stats.Activities) == 0 {
		b.WriteString("No sessions for this period.")
	} else {
		for i, a := range stats.Activities {
			name := a.Name
			if a.Emoji != "" {
				name = a.Emoji + " " + a.Name
			}
			b.WriteString(fmt.Sprintf("%d) %s - %s (%s, %d)\n", i+1, name, formatReportDuration(a.Duration), percentOf(a.Duration, total), a.Sessions))
		}
	}
	m.appendGranularityText(ctx, &b, from, to, activityIDs)
	_, _ = m.bot.Send(tgbotapi.NewMessage(ctx.ChatID, b.String()))
}

func (m *Module) ShowPeriodChartReport(ctx *tgctx.MsgContext, from, to time.Time, activityIDs []int64) {
	stats, err := m.tracksvc.GetPeriodReport(ctx.Ctx, ctx.DBUserID, from, to.Add(24*time.Hour), activityIDs)
	if err != nil {
		_, _ = m.bot.Send(tgbotapi.NewMessage(ctx.ChatID, "⚠️ Failed to build period chart."))
		return
	}
	if len(stats.Activities) == 0 {
		_, _ = m.bot.Send(tgbotapi.NewMessage(ctx.ChatID, "📉 No data for selected period."))
		return
	}
	maxDur := time.Duration(1)
	for _, a := range stats.Activities {
		if a.Duration > maxDur {
			maxDur = a.Duration
		}
	}
	var b strings.Builder
	b.WriteString("📉 Period Chart\n\n")
	b.WriteString(fmt.Sprintf("Range: %s..%s\n\n", from.Format("2006-01-02"), to.Format("2006-01-02")))
	total := stats.TotalTracked
	for _, a := range stats.Activities {
		name := a.Name
		if a.Emoji != "" {
			name = a.Emoji + " " + a.Name
		}
		barLen := int((float64(a.Duration) / float64(maxDur)) * 12)
		if barLen < 1 {
			barLen = 1
		}
		b.WriteString(fmt.Sprintf("%s\n%s %s (%s, %d)\n\n", name, strings.Repeat("█", barLen), formatReportDuration(a.Duration), percentOf(a.Duration, total), a.Sessions))
	}
	m.appendGranularityText(ctx, &b, from, to, activityIDs)
	_, _ = m.bot.Send(tgbotapi.NewMessage(ctx.ChatID, b.String()))
}

func (m *Module) ShowPeriodCalendar(ctx *tgctx.MsgContext, month, from, to time.Time) {
	text := fmt.Sprintf("📅 Pick period days\nFrom: %s\nTo: %s", formatDateOrDash(from), formatDateOrDash(to))
	edit := tgbotapi.NewEditMessageTextAndMarkup(
		ctx.ChatID,
		ctx.MessageID,
		text,
		track.TrackReportPeriodCalendarInlineMenu(month, from, to),
	)
	_, _ = m.bot.Send(edit)
}

func (m *Module) appendGranularityText(ctx *tgctx.MsgContext, b *strings.Builder, from, to time.Time, activityIDs []int64) {
	if len(activityIDs) == 0 {
		return
	}

	granularity := "day"
	labelFmt := "2006-01-02"
	if from.Year() != to.Year() {
		granularity = "month"
		labelFmt = "2006-01"
	} else if from.Year() == to.Year() && from.Month() == to.Month() && from.Day() == to.Day() {
		granularity = "hour"
		labelFmt = "15:00"
	}

	buckets, durs, err := m.tracksvc.GetPeriodBuckets(ctx.Ctx, ctx.DBUserID, from, to.Add(24*time.Hour), activityIDs, granularity)
	if err != nil || len(buckets) == 0 {
		return
	}

	switch granularity {
	case "month":
		b.WriteString("\nBy months:\n")
	case "day":
		b.WriteString("\nBy days:\n")
	case "hour":
		b.WriteString("\nBy hours:\n")
	}

	for i := range buckets {
		b.WriteString(fmt.Sprintf("- %s: %s\n", buckets[i].Format(labelFmt), formatReportDuration(durs[i])))
	}
}

func (m *Module) ShowTodayReport(ctx *tgctx.MsgContext) {
	m.ShowTodayChart(ctx)
}

func (m *Module) ShowTodayReportBySelected(ctx *tgctx.MsgContext) {
	m.ShowTodaySelectActivities(ctx, map[int64]bool{})
}

func (m *Module) ShowTodaySelectActivities(ctx *tgctx.MsgContext, selected map[int64]bool) {
	items, err := m.tracksvc.ListActivities(ctx.Ctx, ctx.DBUserID)
	if err != nil {
		_, _ = m.bot.Send(tgbotapi.NewMessage(ctx.ChatID, "⚠️ Failed to load activities."))
		return
	}
	text := "🧩 Select activities for today chart"
	if ctx.MessageID > 0 {
		edit := tgbotapi.NewEditMessageTextAndMarkup(
			ctx.ChatID,
			ctx.MessageID,
			text,
			track.TrackTodaySelectActivitiesInlineMenu(items, selected),
		)
		_, _ = m.bot.Send(edit)
		return
	}
	msg := tgbotapi.NewMessage(ctx.ChatID, text)
	msg.ReplyMarkup = track.TrackTodaySelectActivitiesInlineMenu(items, selected)
	_, _ = m.bot.Send(msg)
}

func (m *Module) renderTodayReport(ctx *tgctx.MsgContext, stats models.ReportTodayStats, err error, title string) {
	if err != nil {
		log.Error().Err(err).Msg("today report failed")
		if ctx.MessageID > 0 {
			edit := tgbotapi.NewEditMessageText(ctx.ChatID, ctx.MessageID, "⚠️ Failed to load today report.")
			_, _ = m.bot.Send(edit)
		} else {
			_, _ = m.bot.Send(tgbotapi.NewMessage(ctx.ChatID, "⚠️ Failed to load today report."))
		}
		return
	}

	var b strings.Builder
	b.WriteString(title + "\n\n")
	b.WriteString(fmt.Sprintf("Total: %s\n", formatReportDuration(stats.TotalTracked)))
	b.WriteString(fmt.Sprintf("Sessions: %d\n\n", stats.TotalSessions))
	if len(stats.TopActivities) == 0 {
		b.WriteString("Top activities: none yet")
	} else {
		b.WriteString("Top activities:\n")
		for i, item := range stats.TopActivities {
			name := item.Name
			if item.Emoji != "" {
				name = item.Emoji + " " + item.Name
			}
			b.WriteString(fmt.Sprintf("%d) %s - %s (%d)\n", i+1, name, formatReportDuration(item.Duration), item.Sessions))
		}
	}

	if ctx.MessageID > 0 {
		edit := tgbotapi.NewEditMessageTextAndMarkup(
			ctx.ChatID,
			ctx.MessageID,
			b.String(),
			track.TrackReportTodayInlineMenu(),
		)
		_, _ = m.bot.Send(edit)
		return
	}
	msg := tgbotapi.NewMessage(ctx.ChatID, b.String())
	msg.ReplyMarkup = track.TrackReportTodayInlineMenu()
	_, _ = m.bot.Send(msg)
}

func (m *Module) PromptCreateActivity(ctx *tgctx.MsgContext) {
	text := "📌 *Create New Activity*\n\nEnter activity name:"
	msg := tgbotapi.NewMessage(ctx.ChatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = track.TrackActivityManageReplyMenu()

	if _, err := m.bot.Send(msg); err != nil {
		log.Error().Err(err).Msg("send create activity prompt failed")
	}
}

func (m *Module) ProcessCreateActivity(ctx *tgctx.MsgContext) bool {
	name := strings.TrimSpace(ctx.Text)
	if name == "" {
		msg := tgbotapi.NewMessage(ctx.ChatID, "Activity name cannot be empty.")
		_, _ = m.bot.Send(msg)
		return false
	}

	activity, err := m.tracksvc.CreateActivity(ctx.Ctx, ctx.DBUserID, name, "")
	if err != nil {
		if err == models.ErrActivityExists {
			_, _ = m.bot.Send(tgbotapi.NewMessage(ctx.ChatID, "Activity already exists."))
			return false
		}
		log.Error().Err(err).Msg("create activity failed")
		_, _ = m.bot.Send(tgbotapi.NewMessage(ctx.ChatID, "⚠️ Failed to create activity."))
		return false
	}

	confirm := tgbotapi.NewMessage(ctx.ChatID, fmt.Sprintf("Created: %s", activity.Name))
	confirm.ReplyMarkup = track.TrackCreateSuccessInlineMenu()
	_, _ = m.bot.Send(confirm)
	return true
}

func (m *Module) ShowTrackActivitySelectionMenu(ctx *tgctx.MsgContext) {
	items, err := m.tracksvc.ListActivities(ctx.Ctx, ctx.DBUserID)
	if err != nil {
		log.Error().Err(err).Msg("list activities failed")
		_, _ = m.bot.Send(tgbotapi.NewMessage(ctx.ChatID, "⚠️ Failed to load activities."))
		return
	}

	if len(items) == 0 {
		msg := tgbotapi.NewMessage(ctx.ChatID, "No activities yet. Create one first.")
		msg.ReplyMarkup = track.TrackActivityManageReplyMenu()
		_, _ = m.bot.Send(msg)
		return
	}

	selectedCount := 0
	for _, item := range items {
		if item.Selected {
			selectedCount++
		}
	}

	msgReply := tgbotapi.NewMessage(ctx.ChatID, "🗂")
	msgReply.ReplyMarkup = track.TrackActivityManageReplyMenu()
	_, _ = m.bot.Send(msgReply)

	msg := tgbotapi.NewMessage(ctx.ChatID, fmt.Sprintf("📂 Select Activity\n\nSelected: %d of %d", selectedCount, len(items)))
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = track.TrackActivitiesInlineMenu(items)
	_, _ = m.bot.Send(msg)
}

func (m *Module) HandleTrackToggleCallback(ctx *tgctx.MsgContext) {
	payload := strings.TrimPrefix(ctx.Text, "act_toggle_:")
	activityID, err := strconv.ParseInt(payload, 10, 64)
	if err != nil {
		_, _ = m.bot.Send(tgbotapi.NewMessage(ctx.ChatID, "Invalid activity id."))
		return
	}

	if err := m.tracksvc.ToggleSelectedActivity(ctx.Ctx, ctx.DBUserID, activityID); err != nil {
		log.Error().Err(err).Msg("toggle activity failed")
		_, _ = m.bot.Send(tgbotapi.NewMessage(ctx.ChatID, "⚠️ Failed to update activity selection."))
		return
	}

	items, err := m.tracksvc.ListActivities(ctx.Ctx, ctx.DBUserID)
	if err != nil {
		log.Error().Err(err).Msg("reload activities failed")
		_, _ = m.bot.Send(tgbotapi.NewMessage(ctx.ChatID, "⚠️ Failed to refresh activities."))
		return
	}

	selectedCount := 0
	for _, item := range items {
		if item.Selected {
			selectedCount++
		}
	}

	edit := tgbotapi.NewEditMessageTextAndMarkup(
		ctx.ChatID,
		ctx.MessageID,
		fmt.Sprintf("📂 Select Activity\n\nSelected: %d of %d", selectedCount, len(items)),
		track.TrackActivitiesInlineMenu(items),
	)
	edit.ParseMode = "HTML"
	if _, err := m.bot.Send(edit); err != nil {
		log.Error().Err(err).Msg("edit activity list failed")
	}
}

func (m *Module) DeleteSelectedActivities(ctx *tgctx.MsgContext) {
	deleted, err := m.tracksvc.DeleteSelectedActivities(ctx.Ctx, ctx.DBUserID)
	if err != nil {
		log.Error().Err(err).Msg("delete selected activities failed")
		_, _ = m.bot.Send(tgbotapi.NewMessage(ctx.ChatID, "⚠️ Failed to delete selected activities."))
		return
	}

	_, _ = m.bot.Send(tgbotapi.NewMessage(ctx.ChatID, fmt.Sprintf("🗑 Deleted: %d", deleted)))
	m.ShowTrackActivitySelectionMenu(ctx)
}

func (m *Module) ArchiveSelectedActivities(ctx *tgctx.MsgContext) {
	archived, err := m.tracksvc.ArchiveSelectedActivities(ctx.Ctx, ctx.DBUserID)
	if err != nil {
		log.Error().Err(err).Msg("archive selected activities failed")
		_, _ = m.bot.Send(tgbotapi.NewMessage(ctx.ChatID, "⚠️ Failed to archive selected activities."))
		return
	}

	if archived == 0 {
		msg := tgbotapi.NewMessage(ctx.ChatID, "No selected activities to archive.")
		msg.ReplyMarkup = track.TrackArchiveSuccessInlineMenu()
		_, _ = m.bot.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(ctx.ChatID, fmt.Sprintf("📦 Archived: %d", archived))
	msg.ReplyMarkup = track.TrackArchiveSuccessInlineMenu()
	_, _ = m.bot.Send(msg)
}

func (m *Module) ArchiveSelectedActivitiesInPlace(ctx *tgctx.MsgContext) {
	archived, err := m.tracksvc.ArchiveSelectedActivities(ctx.Ctx, ctx.DBUserID)
	if err != nil {
		log.Error().Err(err).Msg("archive selected activities failed")
		edit := tgbotapi.NewEditMessageText(ctx.ChatID, ctx.MessageID, "⚠️ Failed to archive selected activities.")
		_, _ = m.bot.Send(edit)
		return
	}

	if archived == 0 {
		edit := tgbotapi.NewEditMessageTextAndMarkup(
			ctx.ChatID,
			ctx.MessageID,
			"No selected activities to archive.",
			track.TrackArchiveSuccessInlineMenu(),
		)
		_, _ = m.bot.Send(edit)
		return
	}

	edit := tgbotapi.NewEditMessageTextAndMarkup(
		ctx.ChatID,
		ctx.MessageID,
		fmt.Sprintf("📦 Archived: %d", archived),
		track.TrackArchiveSuccessInlineMenu(),
	)
	_, _ = m.bot.Send(edit)
}

func (m *Module) ShowArchiveMenu(ctx *tgctx.MsgContext) {
	m.renderArchiveMenu(ctx, false)
}

func (m *Module) ShowArchiveMenuInPlace(ctx *tgctx.MsgContext) {
	m.renderArchiveMenu(ctx, true)
}

func (m *Module) renderArchiveMenu(ctx *tgctx.MsgContext, edit bool) {
	items, err := m.tracksvc.ListArchivedActivities(ctx.Ctx, ctx.DBUserID)
	if err != nil {
		log.Error().Err(err).Msg("list archive failed")
		if edit && ctx.MessageID > 0 {
			msg := tgbotapi.NewEditMessageText(ctx.ChatID, ctx.MessageID, "⚠️ Failed to load archive.")
			_, _ = m.bot.Send(msg)
		} else {
			_, _ = m.bot.Send(tgbotapi.NewMessage(ctx.ChatID, "⚠️ Failed to load archive."))
		}
		return
	}

	if len(items) == 0 {
		text := "Archive is empty."
		if edit && ctx.MessageID > 0 {
			msg := tgbotapi.NewEditMessageText(ctx.ChatID, ctx.MessageID, text)
			_, _ = m.bot.Send(msg)
		} else {
			_, _ = m.bot.Send(tgbotapi.NewMessage(ctx.ChatID, text))
		}
		return
	}

	text := fmt.Sprintf("🗄 Archive\n\nTotal archived: %d", len(items))
	if edit && ctx.MessageID > 0 {
		msgReply := tgbotapi.NewMessage(ctx.ChatID, "🗄")
		msgReply.ReplyMarkup = track.TrackArchiveReplyMenu()
		_, _ = m.bot.Send(msgReply)

		msg := tgbotapi.NewEditMessageTextAndMarkup(
			ctx.ChatID,
			ctx.MessageID,
			text,
			track.TrackArchiveInlineMenu(items),
		)
		_, _ = m.bot.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(ctx.ChatID, text)
	msg.ReplyMarkup = track.TrackArchiveInlineMenu(items)
	msgReply := tgbotapi.NewMessage(ctx.ChatID, "🗄")
	msgReply.ReplyMarkup = track.TrackArchiveReplyMenu()
	_, _ = m.bot.Send(msgReply)
	_, _ = m.bot.Send(msg)
}

func (m *Module) ShowTrackActivitySelectionMenuInPlace(ctx *tgctx.MsgContext) {
	msgReply := tgbotapi.NewMessage(ctx.ChatID, "🗂")
	msgReply.ReplyMarkup = track.TrackActivityManageReplyMenu()
	_, _ = m.bot.Send(msgReply)

	items, err := m.tracksvc.ListActivities(ctx.Ctx, ctx.DBUserID)
	if err != nil {
		log.Error().Err(err).Msg("list activities failed")
		edit := tgbotapi.NewEditMessageText(ctx.ChatID, ctx.MessageID, "⚠️ Failed to load activities.")
		_, _ = m.bot.Send(edit)
		return
	}

	if len(items) == 0 {
		edit := tgbotapi.NewEditMessageText(ctx.ChatID, ctx.MessageID, "No activities yet. Create one first.")
		_, _ = m.bot.Send(edit)
		return
	}

	selectedCount := 0
	for _, item := range items {
		if item.Selected {
			selectedCount++
		}
	}

	edit := tgbotapi.NewEditMessageTextAndMarkup(
		ctx.ChatID,
		ctx.MessageID,
		fmt.Sprintf("📂 Select Activity\n\nSelected: %d of %d", selectedCount, len(items)),
		track.TrackActivitiesInlineMenu(items),
	)
	edit.ParseMode = "HTML"
	_, _ = m.bot.Send(edit)
}

func (m *Module) RestoreArchivedActivity(ctx *tgctx.MsgContext) {
	idRaw := strings.TrimPrefix(ctx.Text, track.TrackCBArchiveRestore)
	activityID, err := strconv.ParseInt(idRaw, 10, 64)
	if err != nil {
		_, _ = m.bot.Send(tgbotapi.NewMessage(ctx.ChatID, "Invalid activity."))
		return
	}
	activityName := m.findArchivedActivityName(ctx, activityID)

	if err := m.tracksvc.RestoreArchivedActivity(ctx.Ctx, ctx.DBUserID, activityID); err != nil {
		log.Error().Err(err).Msg("restore archived activity failed")
		edit := tgbotapi.NewEditMessageText(ctx.ChatID, ctx.MessageID, "⚠️ Failed to restore activity.")
		_, _ = m.bot.Send(edit)
		return
	}
	_, _ = m.bot.Send(tgbotapi.NewMessage(ctx.ChatID, fmt.Sprintf("♻ Activity restored: %s", activityName)))
	m.ShowArchiveMenuInPlace(ctx)
}

func (m *Module) DeleteArchivedForever(ctx *tgctx.MsgContext) {
	idRaw := strings.TrimPrefix(ctx.Text, track.TrackCBArchiveDelete)
	activityID, err := strconv.ParseInt(idRaw, 10, 64)
	if err != nil {
		_, _ = m.bot.Send(tgbotapi.NewMessage(ctx.ChatID, "Invalid activity."))
		return
	}
	activityName := m.findArchivedActivityName(ctx, activityID)

	if err := m.tracksvc.DeleteArchivedForever(ctx.Ctx, ctx.DBUserID, activityID); err != nil {
		log.Error().Err(err).Msg("delete archived forever failed")
		edit := tgbotapi.NewEditMessageText(ctx.ChatID, ctx.MessageID, "⚠️ Failed to delete activity forever.")
		_, _ = m.bot.Send(edit)
		return
	}
	_, _ = m.bot.Send(tgbotapi.NewMessage(ctx.ChatID, fmt.Sprintf("🗑 Deleted forever: %s", activityName)))
	m.ShowArchiveMenuInPlace(ctx)
}

func (m *Module) findArchivedActivityName(ctx *tgctx.MsgContext, activityID int64) string {
	items, err := m.tracksvc.ListArchivedActivities(ctx.Ctx, ctx.DBUserID)
	if err != nil {
		return fmt.Sprintf("#%d", activityID)
	}
	for _, item := range items {
		if item.ID == activityID {
			if item.Emoji != "" {
				return item.Emoji + " " + item.Name
			}
			return item.Name
		}
	}
	return fmt.Sprintf("#%d", activityID)
}

func (m *Module) ShowTrackTimerMenu(ctx *tgctx.MsgContext) {
	msg := tgbotapi.NewMessage(ctx.ChatID, "Select tracking interval:")
	msg.ReplyMarkup = track.TrackTimerReplyMenu()
	_, _ = m.bot.Send(msg)
}

func (m *Module) ActivateTrackTimer(ctx *tgctx.MsgContext, intervalMin int) {
	if m.testTimerMin > 0 {
		intervalMin = m.testTimerMin
	}

	items, err := m.tracksvc.ListSelectedActivities(ctx.Ctx, ctx.DBUserID)
	if err != nil {
		log.Error().Err(err).Msg("load selected activities failed")
		_, _ = m.bot.Send(tgbotapi.NewMessage(ctx.ChatID, "⚠️ Failed to load selected activities."))
		return
	}
	if len(items) == 0 {
		_, _ = m.bot.Send(tgbotapi.NewMessage(ctx.ChatID, "Select at least one activity before activating timer."))
		return
	}

	if err := m.timersvc.Activate(ctx.Ctx, ctx.DBUserID, intervalMin); err != nil {
		log.Error().Err(err).Msg("activate timer failed")
		_, _ = m.bot.Send(tgbotapi.NewMessage(ctx.ChatID, "⚠️ Failed to activate timer."))
		return
	}

	_, _ = m.bot.Send(tgbotapi.NewMessage(ctx.ChatID, fmt.Sprintf("✅ Timer activated: every %d min", intervalMin)))
	hide := tgbotapi.NewMessage(ctx.ChatID, " ")
	hide.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	_, _ = m.bot.Send(hide)
	m.ShowEntryMenu(ctx)
}

func (m *Module) StopTrackTimer(ctx *tgctx.MsgContext) {
	if err := m.timersvc.Stop(ctx.Ctx, ctx.DBUserID); err != nil {
		log.Error().Err(err).Msg("stop timer failed")
		_, _ = m.bot.Send(tgbotapi.NewMessage(ctx.ChatID, "⚠️ Failed to stop timer."))
		return
	}
	_, _ = m.bot.Send(tgbotapi.NewMessage(ctx.ChatID, "⏹ Timer stopped"))
}

func (m *Module) SendPromptMessage(ctx context.Context, chatID int64, userID int64, intervalMin int) error {
	items, err := m.tracksvc.ListSelectedActivities(ctx, userID)
	if err != nil {
		return err
	}
	if len(items) == 0 {
		return nil
	}

	msg := tgbotapi.NewMessage(chatID, "What are you doing now?")
	msg.ReplyMarkup = track.TrackPromptInlineMenu(items, intervalMin)
	_, err = m.bot.Send(msg)
	return err
}

func (m *Module) RecordPromptAnswer(ctx *tgctx.MsgContext) {
	payload := strings.TrimPrefix(ctx.Text, track.TrackCBPromptActivity)
	parts := strings.Split(payload, ":")
	if len(parts) != 2 {
		_, _ = m.bot.Send(tgbotapi.NewMessage(ctx.ChatID, "Invalid selection payload."))
		return
	}

	activityID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		_, _ = m.bot.Send(tgbotapi.NewMessage(ctx.ChatID, "Invalid activity id."))
		return
	}

	intervalMin, err := strconv.Atoi(parts[1])
	if err != nil || intervalMin <= 0 {
		_, _ = m.bot.Send(tgbotapi.NewMessage(ctx.ChatID, "Invalid interval."))
		return
	}

	if err := m.timersvc.RecordPromptAnswerWithInterval(ctx.Ctx, ctx.DBUserID, activityID, intervalMin); err != nil {
		log.Error().Err(err).Msg("record prompt answer failed")
		_, _ = m.bot.Send(tgbotapi.NewMessage(ctx.ChatID, "⚠️ Failed to save activity."))
		return
	}

	if ctx.MessageID > 0 {
		del := tgbotapi.NewDeleteMessage(ctx.ChatID, ctx.MessageID)
		_, _ = m.bot.Request(del)
	}

	endAt := time.Now()
	startAt := endAt.Add(-time.Duration(intervalMin) * time.Minute)
	activityName := m.findActivityName(ctx, activityID)

	text := fmt.Sprintf(
		"Saved ✅\nActivity: %s\nTime: %s-%s (%d min)",
		activityName,
		startAt.Format("15:04"),
		endAt.Format("15:04"),
		intervalMin,
	)
	_, _ = m.bot.Send(tgbotapi.NewMessage(ctx.ChatID, text))
}

func (m *Module) findActivityName(ctx *tgctx.MsgContext, activityID int64) string {
	items, err := m.tracksvc.ListActivities(ctx.Ctx, ctx.DBUserID)
	if err != nil {
		return fmt.Sprintf("#%d", activityID)
	}
	for _, item := range items {
		if item.ID == activityID {
			if item.Emoji != "" {
				return item.Emoji + " " + item.Name
			}
			return item.Name
		}
	}
	return fmt.Sprintf("#%d", activityID)
}

func formatReportDuration(d time.Duration) string {
	if d < 0 {
		d = 0
	}
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	switch {
	case h > 0 && m > 0:
		return fmt.Sprintf("%dh %dm", h, m)
	case h > 0:
		return fmt.Sprintf("%dh", h)
	default:
		return fmt.Sprintf("%dm", m)
	}
}

func formatDateOrDash(t time.Time) string {
	if t.IsZero() {
		return "—"
	}
	return t.Format("2006-01-02")
}

func percentOf(part, total time.Duration) string {
	if total <= 0 || part <= 0 {
		return "0%"
	}
	p := (float64(part) / float64(total)) * 100.0
	return fmt.Sprintf("%.1f%%", p)
}

func (m *Module) ShowLearningMenu(ctx *tgctx.MsgContext) {
	stats, err := m.learningsvc.GetLearningStats(ctx.Ctx, ctx.UserID)
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
	stats, err := m.subscriptionsvc.GetSubscriptionStats(ctx.Ctx, ctx.UserID)
	if err != nil {
		log.Error().Err(err).Msg("GetSubscriptionStats failed")
		msg := tgbotapi.NewMessage(ctx.ChatID, "⚠️ Failed to load subscription data. Please try again.")
		_, _ = m.bot.Send(msg)
		return
	}

	text := subscription.SubscriptionMenuText(stats)

	msg := tgbotapi.NewMessage(ctx.ChatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = subscription.SubscriptionEntryInlineMenu()

	if _, err := m.bot.Send(msg); err != nil {
		log.Error().Err(err).Msg("send subscription menu failed")
	}
}
