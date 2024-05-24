package model

import "time"

type (
	Balance struct {
		Current   float64 `json:"current"`
		Withdrawn float64 `json:"withdrawn,omitempty"`
	}

	Withdrawal struct {
		Order       string    `json:"order"`
		Sum         float64   `json:"sum,omitempty"`
		ProcessedAt time.Time `json:"processed_at,omitempty"`
	}

	Withdrawals []Withdrawal
)
