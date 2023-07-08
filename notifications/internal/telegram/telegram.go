package telegram

import (
	"context"
	"route256/libs/tracer"

	"github.com/go-telegram/bot"
	"github.com/opentracing/opentracing-go"
)

type telegramBot struct {
	bot              *bot.Bot
	reciever_chat_id string
}

func (t *telegramBot) SendMessage(ctx context.Context, text string) error {

	span, ctx := opentracing.StartSpanFromContext(ctx, "telegram/SendMessage")
	defer span.Finish()

	_, err := t.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: t.reciever_chat_id,
		Text:   text})

	if err != nil {
		return tracer.MarkSpanWithError(ctx, err)
	}

	return nil

}

func New(ctx context.Context, token, recieverChatId string) (*telegramBot, error) {

	span, ctx := opentracing.StartSpanFromContext(ctx, "telegram/SendMessage")
	defer span.Finish()

	myBot, err := bot.New(token)

	if err != nil {
		return nil, tracer.MarkSpanWithError(ctx, err)
	}

	return &telegramBot{myBot, recieverChatId}, nil
}
