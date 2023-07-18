package schema

import "time"

type StatusChangeMessage struct {
	OldStatus string    `db:"oldstatus"`
	NewStatus string    `db:"newstatus"`
	OrderID   int64     `db:"orderid"`
	DT        time.Time `db:"dt"`
	UserID    int64     `db:"userid"`
}
