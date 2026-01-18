package hendlers

import (
	router "tracker-bot/internal/hendlers"
	tgclient "tracker-bot/internal/utils/tgcient"
	"tracker-bot/internal/utils/tgctx"

	"github.com/rs/zerolog/log"
)

type ReplyModule struct {
	bot          tgclient.TgBotAPI
	track        *router.Module
	subscription *router.Module
	entry        *router.Module
	profile      *router.Module
	learning     *router.Module
}

func New(bot tgclient.TgBotAPI, track *router.Module, subscription *router.Module, entry *router.Module, profile *router.Module, learning *router.Module) *ReplyModule {
	return &ReplyModule{
		bot:          bot,
		track:        track,
		subscription: subscription,
		entry:        entry,
		profile:      profile,
		learning:     learning,
	}
}
func (r *ReplyModule) HandleReplyButtons(ctx *tgctx.MsgContext) bool {
	replyButtons := map[string]func(*tgctx.MsgContext){
		"📈Track": r.handleShowTrackingMenu,
	}
	if handler, ok := replyButtons[ctx.Text]; ok {
		handler(ctx)
		return true
	}
	log.Warn().Msgf("Unknown reply button: %s", ctx.Text)
	return false
}

// reply button

func (r *ReplyModule) handleShowTrackingMenu(ctx *tgctx.MsgContext) {
	r.track.ShowTrackingMenu(ctx)
}
