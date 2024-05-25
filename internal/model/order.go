package model

import "time"

type (
	OrderStatus int

	Order struct {
		Number     string    `json:"number"`
		Status     string    `json:"status"`
		Accrual    float64   `json:"accrual,omitempty"`
		UploadedAt time.Time `json:"uploaded_at,omitempty"`
	}

	Orders []Order
)

const (
	NEW OrderStatus = iota + 1
	PROCESSING
	INVALID
	PROCESSED
)

func (os OrderStatus) String() string {
	return [...]string{"NEW", "PROCESSING", "INVALID", "PROCESSED"}[os-1]
}

func (os OrderStatus) Index() int {
	return int(os)
}
