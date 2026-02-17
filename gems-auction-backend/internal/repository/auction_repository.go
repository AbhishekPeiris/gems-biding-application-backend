package repository

import (
	"context"
	"time"

	"github.com/boswin/gems-auction-backend/config"
	"github.com/boswin/gems-auction-backend/internal/domain"
)

type AuctionRepository struct{}

func NewAuctionRepository() *AuctionRepository {
	return &AuctionRepository{}
}

func (r *AuctionRepository) Create(a *domain.Auction) error {
	query := `
		INSERT INTO auctions (gem_id,start_price,current_price,min_increment,start_time,end_time,status,created_at,updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		RETURNING id
	`

	now := time.Now()

	return config.DB.QueryRow(context.Background(), query,
		a.GemID,
		a.StartPrice,
		a.CurrentPrice,
		a.MinIncrement,
		a.StartTime,
		a.EndTime,
		a.Status,
		now,
		now,
	).Scan(&a.ID)
}

func (r *AuctionRepository) UpdateCurrentPrice(id int64, price float64) error {
	query := `UPDATE auctions SET current_price=$1, updated_at=$2 WHERE id=$3`
	_, err := config.DB.Exec(context.Background(), query, price, time.Now(), id)
	return err
}
