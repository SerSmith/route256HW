package domain

import (
	"context"
	"time"
)

const ServiceName = "Notifications"

//go:generate mockery --name Repository
type Repository interface {
	ReadNotifications(ctx context.Context, req NotificationHistoryRequest) ([]NotificationMem, error)
	WriteNotification(ctx context.Context, message StatusChangeMessage, DT time.Time) error
}

//go:generate mockery --name CashDB
type CashDB interface {
	Get(ctx context.Context, req NotificationHistoryRequest) ([]NotificationMem, bool, error)
	Set(ctx context.Context, req NotificationHistoryRequest, value []NotificationMem) error
}

type Messenger interface {
	SendMessage(ctx context.Context, text string) error
}

type Receiver interface {
	Subscribe(ctx context.Context, text string) error
}

type Model struct {
	Messenger  Messenger
	Repository Repository
	CashDB     CashDB
}

func New(ctx context.Context, Messenger Messenger, Repository Repository, CashDB CashDB) *Model {
	return &Model{Messenger, Repository, CashDB}
}

type StatusChangeMessage struct {
	OldStatus string
	NewStatus string
	OrderID   int64
	UserID    int64
}

type NotificationMem struct {
	ChangeStatus *StatusChangeMessage
	DT           time.Time
}

type NotificationHistoryRequest struct {
	UserID   int64
	DateFrom time.Time
	DateTo   time.Time
}
