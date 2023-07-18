package service

import (
	"context"
	// "errors"
	"log"
	"route256/notifications/internal/domain"
	"route256/notifications/internal/domain/mocks"
	"route256/notifications/internal/handlers/http/getNotificationHistory"
	"route256/notifications/internal/handlers/http/schema"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_getNotificationHistory(t *testing.T) {
	t.Parallel()

	t.Run("success, get from repo", func(t *testing.T) {
		t.Parallel()

		const (
			userID   = int64(1)
			DateFrom = "2023-07-12"
			DateTo   = "2023-07-14"
		)

		date, err := time.Parse("2006-01-02", "2023-11-30")

		require.NoError(t, err)

		DateFromTime, err := time.Parse("2006-01-02", DateFrom)

		require.NoError(t, err)

		DateToTime, err := time.Parse("2006-01-02", DateTo)

		require.NoError(t, err)

		var (
			req = schema.Request{UserID: userID,
				DateFrom: DateFrom,
				DateTo:   DateTo}

			csm1 = domain.StatusChangeMessage{OldStatus: "Oldstatus1",
				NewStatus: "NewStatus1",
				OrderID:   1,
				UserID:    1}

			ResponseEl1 = schema.ResponseEl{OldStatus: csm1.OldStatus,
				NewStatus: csm1.NewStatus,
				OrderID:   csm1.OrderID,
				UserID:    csm1.UserID,
				DT:        date}

			correctAnswer = schema.Response{Data: []schema.ResponseEl{ResponseEl1}}
		)

		var (
			reqDB = domain.NotificationHistoryRequest{UserID: userID,
				DateFrom: DateFromTime,
				DateTo:   DateToTime}

			outRepo = domain.NotificationMem{ChangeStatus: &csm1,
				DT: date}
		)

		repositoryMock := mocks.NewRepository(t)
		repositoryMock.On("ReadNotifications", mock.Anything, reqDB).Return([]domain.NotificationMem{outRepo}, nil).Once()

		CashDBMock := mocks.NewCashDB(t)
		CashDBMock.On("Get", mock.Anything, reqDB).Return(nil, false, nil).Once()
		CashDBMock.On("Set", mock.Anything, mock.Anything, mock.Anything).Return(nil)

		m := domain.New(context.Background(), nil, repositoryMock, CashDBMock)

		getNotificationHistory := getNotificationHistory.Handler{Model: m}

		// Act

		ans, err := getNotificationHistory.Handle(context.Background(), req)

		// Assert
		require.NoError(t, err)
		require.ElementsMatch(t, ans.Data, correctAnswer.Data)
	})

	t.Run("success, get from cash", func(t *testing.T) {
		t.Parallel()

		const (
			userID   = int64(1)
			DateFrom = "2023-07-12"
			DateTo   = "2023-07-14"
		)

		date, err := time.Parse("2006-01-02", "2023-11-30")

		log.Print(date)

		require.NoError(t, err)

		DateFromTime, err := time.Parse("2006-01-02", DateFrom)

		require.NoError(t, err)

		DateToTime, err := time.Parse("2006-01-02", DateTo)

		require.NoError(t, err)

		var (
			req = schema.Request{UserID: userID,
				DateFrom: DateFrom,
				DateTo:   DateTo}

			csm2 = domain.StatusChangeMessage{OldStatus: "status3",
				NewStatus: "status2",
				OrderID:   1,
				UserID:    1}

			ResponseEl2 = schema.ResponseEl{OldStatus: csm2.OldStatus,
				NewStatus: csm2.NewStatus,
				OrderID:   csm2.OrderID,
				UserID:    csm2.UserID,
				DT:        date}

			correctAnswer = schema.Response{Data: []schema.ResponseEl{ResponseEl2}}
		)

		var (
			reqDB = domain.NotificationHistoryRequest{UserID: userID,
				DateFrom: DateFromTime,
				DateTo:   DateToTime}
		)

		repositoryMock := mocks.NewRepository(t)

		var (
			outCash = domain.NotificationMem{ChangeStatus: &csm2,
				DT: date}
		)

		CashDBMock := mocks.NewCashDB(t)
		CashDBMock.On("Get", mock.Anything, reqDB).Return([]domain.NotificationMem{outCash}, true, nil)

		m := domain.New(context.Background(), nil, repositoryMock, CashDBMock)

		getNotificationHistory := getNotificationHistory.Handler{Model: m}

		// Act

		ans, err := getNotificationHistory.Handle(context.Background(), req)

		// Assert
		require.NoError(t, err)
		require.ElementsMatch(t, ans.Data, correctAnswer.Data)
	})

	// t.Run("suc

}
