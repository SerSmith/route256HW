package domain

import (
	"context"
)

type Messenger interface {
	SendMessage(ctx context.Context, text string) error
}

type Receiver interface {
	Subscribe(ctx context.Context, text string) error
}

type Model struct {
	Messenger Messenger
}

func New(ctx context.Context, Messenger Messenger) *Model {
	return &Model{Messenger}
}

type StatusChangeMessage struct {
	OldStatus string
	NewStatus string
	OrderID   int64
}
