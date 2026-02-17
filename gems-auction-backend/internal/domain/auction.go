package domain

import "time"

type AuctionStatus string

const (
	AuctionScheduled AuctionStatus = "SCHEDULED"
	AuctionLive      AuctionStatus = "LIVE"
	AuctionEnded     AuctionStatus = "ENDED"
)

type Auction struct {
	ID            int64         `json:"id"`
	GemID         int64         `json:"gem_id"`
	StartPrice    float64       `json:"start_price"`
	CurrentPrice  float64       `json:"current_price"`
	MinIncrement  float64       `json:"min_increment"`
	StartTime     time.Time     `json:"start_time"`
	EndTime       time.Time     `json:"end_time"`
	Status        AuctionStatus `json:"status"`
	WinnerID      *int64        `json:"winner_id,omitempty"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}
