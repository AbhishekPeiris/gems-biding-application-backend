package repository

import (
	"context"
	"time"

	"github.com/boswin/gems-auction-backend/config"
	"github.com/boswin/gems-auction-backend/internal/domain"
)

type GemRepository struct{}

func NewGemRepository() *GemRepository {
	return &GemRepository{}
}

func (r *GemRepository) Create(gem *domain.Gem) error {
	query := `
		INSERT INTO gems (seller_id,name,description,carat,color,clarity,origin,certificate,image_url,status,created_at,updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
		RETURNING id
	`

	now := time.Now()

	return config.DB.QueryRow(context.Background(), query,
		gem.SellerID,
		gem.Name,
		gem.Description,
		gem.Carat,
		gem.Color,
		gem.Clarity,
		gem.Origin,
		gem.Certificate,
		gem.ImageURL,
		gem.Status,
		now,
		now,
	).Scan(&gem.ID)
}

func (r *GemRepository) GetByID(id int64) (*domain.Gem, error) {
	query := `SELECT id,seller_id,name,description,carat,color,clarity,origin,certificate,image_url,status,created_at,updated_at FROM gems WHERE id=$1`

	var gem domain.Gem

	err := config.DB.QueryRow(context.Background(), query, id).Scan(
		&gem.ID,
		&gem.SellerID,
		&gem.Name,
		&gem.Description,
		&gem.Carat,
		&gem.Color,
		&gem.Clarity,
		&gem.Origin,
		&gem.Certificate,
		&gem.ImageURL,
		&gem.Status,
		&gem.CreatedAt,
		&gem.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &gem, nil
}
