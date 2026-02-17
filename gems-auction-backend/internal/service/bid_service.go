package service

import (
	"context"
	"errors"
	"time"

	"github.com/boswin/gems-auction-backend/config"
	"github.com/boswin/gems-auction-backend/internal/domain"
	"github.com/boswin/gems-auction-backend/internal/repository"
	"github.com/jackc/pgx/v5"
)

type AuctionEventBroadcaster interface {
	BroadcastToAuction(auctionID int64, eventType string, payload any)
}

type BidService struct {
	bidRepo     *repository.BidRepository
	auctionRepo *repository.AuctionRepository
	broadcast   AuctionEventBroadcaster // can be nil for now
}

func NewBidService(
	bidRepo *repository.BidRepository,
	auctionRepo *repository.AuctionRepository,
	broadcast AuctionEventBroadcaster,
) *BidService {
	return &BidService{bidRepo: bidRepo, auctionRepo: auctionRepo, broadcast: broadcast}
}

type PlaceBidRequest struct {
	AuctionID int64   `json:"auction_id"`
	UserID    int64   `json:"user_id"`
	Amount    float64 `json:"amount"`
}

type BidPlacedEvent struct {
	AuctionID   int64     `json:"auction_id"`
	UserID      int64     `json:"user_id"`
	Amount      float64   `json:"amount"`
	PlacedAt    time.Time `json:"placed_at"`
	NewHighBid  float64   `json:"new_high_bid"`
}

func (s *BidService) PlaceBid(req PlaceBidRequest) (*domain.Bid, error) {
	if req.AuctionID <= 0 || req.UserID <= 0 {
		return nil, errors.New("auction_id and user_id required")
	}
	if req.Amount <= 0 {
		return nil, errors.New("amount must be > 0")
	}

	ctx := context.Background()
	tx, err := config.DB.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Lock auction row to avoid race conditions (two users bidding same time)
	var (
		currentPrice float64
		minInc       float64
		status       domain.AuctionStatus
		endTime      time.Time
	)

	q := `SELECT current_price, min_increment, status, end_time
	      FROM auctions WHERE id=$1 FOR UPDATE`

	if err := tx.QueryRow(ctx, q, req.AuctionID).Scan(&currentPrice, &minInc, &status, &endTime); err != nil {
		return nil, err
	}

	if status != domain.AuctionLive {
		return nil, errors.New("auction is not live")
	}
	if time.Now().After(endTime) {
		return nil, errors.New("auction ended")
	}

	minAllowed := currentPrice + minInc
	if req.Amount < minAllowed {
		return nil, errors.New("bid too low (must be at least current_price + min_increment)")
	}

	// Insert bid
	var bidID int64
	ins := `INSERT INTO bids (auction_id, user_id, amount, created_at)
	        VALUES ($1,$2,$3,$4) RETURNING id`
	now := time.Now()

	if err := tx.QueryRow(ctx, ins, req.AuctionID, req.UserID, req.Amount, now).Scan(&bidID); err != nil {
		return nil, err
	}

	// Update auction current_price
	up := `UPDATE auctions SET current_price=$1, updated_at=$2 WHERE id=$3`
	if _, err := tx.Exec(ctx, up, req.Amount, now, req.AuctionID); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	b := &domain.Bid{
		ID:        bidID,
		AuctionID: req.AuctionID,
		UserID:    req.UserID,
		Amount:    req.Amount,
		CreatedAt: now,
	}

	// Broadcast event to websocket clients (optional)
	if s.broadcast != nil {
		s.broadcast.BroadcastToAuction(req.AuctionID, "BID_PLACED", BidPlacedEvent{
			AuctionID:  req.AuctionID,
			UserID:     req.UserID,
			Amount:     req.Amount,
			PlacedAt:   now,
			NewHighBid: req.Amount,
		})
	}

	return b, nil
}
