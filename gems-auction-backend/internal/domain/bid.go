package domain

import "time"

type Bid struct {
	ID        int64     `json:"id"`
	AuctionID int64     `json:"auction_id"`
	UserID    int64     `json:"user_id"`
	Amount    float64   `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}
