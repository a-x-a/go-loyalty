package accrualclient

type AccrualStatus int

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

type AccrualOrder struct {
	Order   int64   `json:"order,string"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}

func (o AccrualOrder) IsValid() bool {
	for _, s := range statuses() {
		if o.Status == s {
			return true
		}
	}

	return false
}
