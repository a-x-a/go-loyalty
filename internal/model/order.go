package model

import "time"

type (
	Order struct {
		Number      string      `json:"number"`
		Status      string      `json:"status"`
		Accrual     float64     `json:"accrual,omitempty"`
		UploadedAt  time.Ticker `json:"uploaded_at,omitempty"`
		ProcessedAt time.Time   `json:"processed_at,omitempty"`
		UserID      int64       `json:"-"`
	}

	Orders []Order

	Withdrawals struct {
		Number      string    `json:"order"`
		Accrual     float64   `json:"sum,omitempty"`
		ProcessedAt time.Time `json:"processed_at,omitempty"`
	}
)
