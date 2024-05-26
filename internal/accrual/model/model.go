package model

type (
	AccrualOrder struct {
		UID     int64   `json:"user_id" db:"user_id"`
		Order   string  `json:"order" db:"number"`
		Status  string  `json:"status" db:"status"`
		Accrual float64 `json:"accrual" db:"accrual"`
	}

	AccrualOrders []AccrualOrder

	AccrualStatus int
)

const (
	REGISTERED AccrualStatus = iota + 1
	PROCESSING
	INVALID
	PROCESSED
)

func statuses() [4]string {
	return [4]string{"REGISTERED", "PROCESSING", "INVALID", "PROCESSED"}
}

func (s AccrualStatus) String() string {
	return statuses()[s-1]
}

func (s AccrualStatus) Index() int {
	return int(s)
}

func (o AccrualOrder) IsValid() bool {
	for _, s := range statuses() {
		if o.Status == s {
			return true
		}
	}

	return false
}

func (o AccrualOrder) GetStatusIndex() AccrualStatus {
	for i, s := range statuses() {
		if o.Status == s {
			return AccrualStatus(i)
		}
	}

	return AccrualStatus(-1)
}
