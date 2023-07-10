package telegram

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
)

type telegramBot struct {
	bot              *bot.Bot
	reciever_chat_id string
}

func (t *telegramBot) SendMessage(ctx context.Context, text string) error {
	_, err := t.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: t.reciever_chat_id,
		Text:   text})

	if err != nil {
		return fmt.Errorf("t.bot.SendMessage err: %w", err)
	}

	return nil

}

func New(ctx context.Context, token, recieverChatId string) (*telegramBot, error) {

	myBot, err := bot.New(token)

	if err != nil {
		return nil, fmt.Errorf("bot.New: %w", err)
	}

	return &telegramBot{myBot, recieverChatId}, nil
}
