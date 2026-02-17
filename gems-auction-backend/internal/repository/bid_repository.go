package repository

import (
	"context"
	"time"

	"github.com/boswin/gems-auction-backend/config"
	"github.com/boswin/gems-auction-backend/internal/domain"
)

type BidRepository struct{}

func NewBidRepository() *BidRepository {
	return &BidRepository{}
}

func (r *BidRepository) Create(b *domain.Bid) error {
	query := `
		INSERT INTO bids (auction_id,user_id,amount,created_at)
		VALUES ($1,$2,$3,$4)
		RETURNING id
	`

	return config.DB.QueryRow(context.Background(), query,
		b.AuctionID,
		b.UserID,
		b.Amount,
		time.Now(),
	).Scan(&b.ID)
}

func (r *BidRepository) GetHighestBid(auctionID int64) (*domain.Bid, error) {
	query := `
		SELECT id,auction_id,user_id,amount,created_at
		FROM bids
		WHERE auction_id=$1
		ORDER BY amount DESC
		LIMIT 1
	`

	var bid domain.Bid

	err := config.DB.QueryRow(context.Background(), query, auctionID).Scan(
		&bid.ID,
		&bid.AuctionID,
		&bid.UserID,
		&bid.Amount,
		&bid.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &bid, nil
}
