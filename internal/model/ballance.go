package model

import "time"

type (
	Balance struct {
		Current   float64 `json:"current,omitempty" example:"500.5"`
		Withdrawn float64 `json:"withdrawn,omitempty" example:"42"`
	} //	@Name Balance

	Withdrawal struct {
		Order       string    `json:"order" example:"2377225624"`
		Sum         float64   `json:"sum,omitempty" example:"500"`
		ProcessedAt time.Time `json:"processed_at,omitempty" example:"2020-12-09T16:09:57+03:00"`
	} //	@Name Withdrawal

	Withdrawals []Withdrawal // @Name Withdrawals
)
