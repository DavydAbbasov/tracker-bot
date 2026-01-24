package tgclient

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TgBotAPI interface {
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
	GetUpdatesChan(config tgbotapi.UpdateConfig) tgbotapi.UpdatesChannel
	StopReceivingUpdates()
	Request(c tgbotapi.Chattable) (*tgbotapi.APIResponse, error)
}

type HandlerFunc func(tgUpdate tgbotapi.Update, c *Client)

func (f HandlerFunc) RunFunc(tgUpdate tgbotapi.Update, c *Client) {
	f(tgUpdate, c)
}

type Client struct {
	client                *tgbotapi.BotAPI
	handlerProcessingFunc HandlerFunc
}

func New(token string) (*Client, error) {
	if token == "" {
		return nil, fmt.Errorf("telegram token is empty")
	}

	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	return &Client{client: api}, nil
}
