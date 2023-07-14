package schema

import (
	"time"
)

type ResponseEl struct {
	OldStatus string
	NewStatus string
	OrderID   int64
	UserID    int64
	DT        time.Time
}

type Response struct {
	Data []ResponseEl
}

func (r Request) Validate() error {
	return nil
}

type Request struct {
	UserID   int64  `json:"UserID"`
	DateFrom string `json:"DateFrom"`
	DateTo   string `json:"DateTo"`
}
