package model

type Balance struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn,omitempty"`
}
