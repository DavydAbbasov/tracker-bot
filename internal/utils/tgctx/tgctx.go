package tgctx

import "context"

// MsgContext is Telegram update context for message-based handlers.
type MsgContext struct {
	Ctx context.Context

	ChatID   int64
	UserID   int64
	DBUserID int64

	Text      string
	MessageID int
}
