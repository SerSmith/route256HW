package schema

import (
	"time"
)

type StatusChangeMessage struct {
	OldStatus string    `json:"OldStatus"`
	NewStatus string    `json:"NewStatus"`
	OrderID   int64     `json:"OrderID"`
	DT        time.Time `json:"DT"`
	UserID    int64     `json:"UserID"`
}

type NotificationHistoryRequest struct {
	UserID   int64     `json:"OldStatus"`
	DateFrom time.Time `json:"DateFrom"`
	DateTo   time.Time `json:"DateTo"`
}
