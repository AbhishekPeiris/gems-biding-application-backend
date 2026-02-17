package domain

import "time"

type PaymentStatus string

const (
	PaymentPending   PaymentStatus = "PENDING"
	PaymentCompleted PaymentStatus = "COMPLETED"
	PaymentFailed    PaymentStatus = "FAILED"
)

type Payment struct {
	ID         int64         `json:"id"`
	AuctionID  int64         `json:"auction_id"`
	UserID     int64         `json:"user_id"`
	Amount     float64       `json:"amount"`
	Status     PaymentStatus `json:"status"`
	Reference  string        `json:"reference"`
	CreatedAt  time.Time     `json:"created_at"`
	UpdatedAt  time.Time     `json:"updated_at"`
}
