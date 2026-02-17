package service

import (
	"context"
	"errors"
	"time"

	"github.com/boswin/gems-auction-backend/config"
	"github.com/boswin/gems-auction-backend/internal/domain"
)

type PaymentService struct{}

func NewPaymentService() *PaymentService {
	return &PaymentService{}
}

type CreatePaymentRequest struct {
	AuctionID int64   `json:"auction_id"`
	UserID    int64   `json:"user_id"`
	Amount    float64 `json:"amount"`
	Reference string  `json:"reference"`
}

// Placeholder flow: create PENDING payment record
func (s *PaymentService) CreatePending(req CreatePaymentRequest) (*domain.Payment, error) {
	if req.AuctionID <= 0 || req.UserID <= 0 {
		return nil, errors.New("auction_id and user_id required")
	}
	if req.Amount <= 0 {
		return nil, errors.New("amount must be > 0")
	}

	now := time.Now()
	p := &domain.Payment{
		AuctionID: req.AuctionID,
		UserID:    req.UserID,
		Amount:    req.Amount,
		Status:    domain.PaymentPending,
		Reference: req.Reference,
		CreatedAt: now,
		UpdatedAt: now,
	}

	q := `INSERT INTO payments (auction_id, user_id, amount, status, reference, created_at, updated_at)
	      VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id`

	if err := config.DB.QueryRow(context.Background(), q,
		p.AuctionID, p.UserID, p.Amount, p.Status, p.Reference, p.CreatedAt, p.UpdatedAt,
	).Scan(&p.ID); err != nil {
		return nil, err
	}

	return p, nil
}

func (s *PaymentService) MarkCompleted(paymentID int64) error {
	if paymentID <= 0 {
		return errors.New("invalid payment id")
	}
	q := `UPDATE payments SET status=$1, updated_at=$2 WHERE id=$3`
	_, err := config.DB.Exec(context.Background(), q, domain.PaymentCompleted, time.Now(), paymentID)
	return err
}

func (s *PaymentService) MarkFailed(paymentID int64) error {
	if paymentID <= 0 {
		return errors.New("invalid payment id")
	}
	q := `UPDATE payments SET status=$1, updated_at=$2 WHERE id=$3`
	_, err := config.DB.Exec(context.Background(), q, domain.PaymentFailed, time.Now(), paymentID)
	return err
}
