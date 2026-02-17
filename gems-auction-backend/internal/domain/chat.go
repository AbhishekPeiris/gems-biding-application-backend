package domain

import "time"

type ChatMessage struct {
	ID        int64     `json:"id"`
	AuctionID int64     `json:"auction_id"`
	UserID    int64     `json:"user_id"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}
