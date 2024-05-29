package model

import "time"

type (
	OrderStatus int

	Order struct {
		Number     string    `json:"number" example:"9278923470"`
		Status     string    `json:"status" enums:"NEW,PROCESSING,INVALID,PROCESSED" example:"PROCESSED"`
		Accrual    float64   `json:"accrual,omitempty" example:"500"`
		UploadedAt time.Time `json:"uploaded_at,omitempty" example:"2020-12-10T15:15:45+03:00"`
	} //	@Name Order

	Orders []Order //	@Name Orders
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
