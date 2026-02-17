package websocket

import "time"

// Event is what we send to frontend via WebSocket
type Event struct {
	Type      string    `json:"type"`
	AuctionID int64     `json:"auction_id"`
	Payload   any       `json:"payload"`
	Timestamp time.Time `json:"timestamp"`
}
