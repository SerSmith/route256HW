package postgres

import (
	"context"
	"fmt"
	"route256/libs/logger"
	"route256/libs/tx"
	"route256/notifications/internal/domain"
	"route256/notifications/internal/repository/schema"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

type Repository struct {
	provider tx.DBProvider
}

func New(provider tx.DBProvider) *Repository {
	return &Repository{provider: provider}
}

const (
	tableNotificationsSent = "notifications_sent"
)

func (r *Repository) ReadNotifications(ctx context.Context, rec domain.NotificationHistoryRequest) ([]domain.NotificationMem, error) {
	db := r.provider.GetDB(ctx)

	query := psql.Select("oldstatus", "newstatus", "orderid", "dt", "userid").
		From(tableNotificationsSent).
		Where(sq.LtOrEq{"dt": rec.DateTo}).
		Where(sq.GtOrEq{"dt": rec.DateFrom}).
		Where(sq.Eq{"userid": rec.UserID})

	rawSQL, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query for ReadNotification: %s", err)
	}

	var StatusMessages []schema.StatusChangeMessage
	err = pgxscan.Select(ctx, db, &StatusMessages, rawSQL, args...)

	if err != nil {
		return nil, fmt.Errorf("exec Select ReadNotification: %w", err)
	}

	out := make([]domain.NotificationMem, 0, len(StatusMessages))

	for _, message := range StatusMessages {

		mes := domain.StatusChangeMessage{
			OldStatus: message.OldStatus,
			NewStatus: message.NewStatus,
			OrderID:   message.OrderID,
			UserID:    message.UserID}

		outMessage := domain.NotificationMem{
			ChangeStatus: &mes,
			DT:           message.DT}

		out = append(out, outMessage)

	}

	return out, nil
}

func (r *Repository) WriteNotification(ctx context.Context, message domain.StatusChangeMessage, DT time.Time) error {
	db := r.provider.GetDB(ctx)

	query := psql.Insert(tableNotificationsSent).
		Columns("oldstatus", "newstatus", "orderid", "dt", "userid").
		Values(message.OldStatus, message.NewStatus, message.OrderID, DT, message.UserID)

	rawSQL, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("build query for WriteNotification: %s", err)
	}

	_, err = db.Exec(ctx, rawSQL, args...)
	if err != nil {
		return fmt.Errorf("exec WriteNotification: %w", err)
	}

	return nil

}
