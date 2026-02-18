package service

import (
	"context"
	"errors"
	"time"

	"github.com/boswin/gems-auction-backend/config"
	"github.com/boswin/gems-auction-backend/internal/domain"
	"github.com/boswin/gems-auction-backend/internal/repository"
)

type AuctionService struct {
	auctionRepo *repository.AuctionRepository
}

func NewAuctionService(auctionRepo *repository.AuctionRepository) *AuctionService {
	return &AuctionService{auctionRepo: auctionRepo}
}

type CreateAuctionRequest struct {
	GemID        int64     `json:"gem_id"`
	StartPrice   float64   `json:"start_price"`
	MinIncrement float64   `json:"min_increment"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
}

func (s *AuctionService) Create(req CreateAuctionRequest) (*domain.Auction, error) {
	if req.GemID <= 0 {
		return nil, errors.New("gem_id required")
	}
	if req.StartPrice <= 0 {
		return nil, errors.New("start_price must be > 0")
	}
	if req.MinIncrement <= 0 {
		return nil, errors.New("min_increment must be > 0")
	}
	if req.EndTime.Before(req.StartTime) || req.EndTime.Equal(req.StartTime) {
		return nil, errors.New("end_time must be after start_time")
	}

	a := &domain.Auction{
		GemID:        req.GemID,
		StartPrice:   req.StartPrice,
		CurrentPrice: req.StartPrice,
		MinIncrement: req.MinIncrement,
		StartTime:    req.StartTime,
		EndTime:      req.EndTime,
		Status:       domain.AuctionScheduled,
	}

	if err := s.auctionRepo.Create(a); err != nil {
		return nil, err
	}

	return a, nil
}

// StartAuction sets status = LIVE (simple helper)
func (s *AuctionService) StartAuction(auctionID int64) error {
	if auctionID <= 0 {
		return errors.New("invalid auction id")
	}

	q := `UPDATE auctions SET status=$1, updated_at=$2 WHERE id=$3`
	_, err := config.DB.Exec(context.Background(), q, domain.AuctionLive, time.Now(), auctionID)
	return err
}

// EndAuction sets status = ENDED and optional winner_id
func (s *AuctionService) EndAuction(auctionID int64, winnerID *int64) error {
	if auctionID <= 0 {
		return errors.New("invalid auction id")
	}

	q := `UPDATE auctions SET status=$1, winner_id=$2, updated_at=$3 WHERE id=$4`
	_, err := config.DB.Exec(context.Background(), q, domain.AuctionEnded, winnerID, time.Now(), auctionID)
	return err
}

func (s *AuctionService) GetByID(auctionID int64) (*domain.Auction, error) {
	if auctionID <= 0 {
		return nil, errors.New("invalid auction id")
	}

	var a domain.Auction
	q := `SELECT id, gem_id, start_price, current_price, min_increment, start_time, end_time, status, winner_id, created_at, updated_at
	      FROM auctions WHERE id=$1`

	err := config.DB.QueryRow(context.Background(), q, auctionID).Scan(
		&a.ID,
		&a.GemID,
		&a.StartPrice,
		&a.CurrentPrice,
		&a.MinIncrement,
		&a.StartTime,
		&a.EndTime,
		&a.Status,
		&a.WinnerID,
		&a.CreatedAt,
		&a.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (s *AuctionService) GetAllAuctions() ([]domain.Auction, error) {
	return s.auctionRepo.GetAll()
}
