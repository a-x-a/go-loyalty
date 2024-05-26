package model

type (
	AccrualOrder struct {
		UID     int64   `db:"user_id" json:"user_id"`
		Order   string  `db:"number" json:"order"`
		Status  string  `db:"status" json:"status"`
		Accrual float64 `db:"accrual" json:"accrual"`
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
