package service

import (
	"errors"

	"github.com/boswin/gems-auction-backend/internal/domain"
	"github.com/boswin/gems-auction-backend/internal/repository"
)

type GemService struct {
	gemRepo *repository.GemRepository
}

func NewGemService(gemRepo *repository.GemRepository) *GemService {
	return &GemService{gemRepo: gemRepo}
}

type CreateGemRequest struct {
	SellerID    int64   `json:"seller_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Carat       float64 `json:"carat"`
	Color       string  `json:"color"`
	Clarity     string  `json:"clarity"`
	Origin      string  `json:"origin"`
	Certificate string  `json:"certificate"`
	ImageURL    string  `json:"image_url"`
}

func (s *GemService) Create(req CreateGemRequest) (*domain.Gem, error) {
	if req.SellerID == 0 {
		return nil, errors.New("seller_id required")
	}
	if req.Name == "" {
		return nil, errors.New("name required")
	}
	if req.Carat <= 0 {
		return nil, errors.New("carat must be > 0")
	}

	g := &domain.Gem{
		SellerID:    req.SellerID,
		Name:        req.Name,
		Description: req.Description,
		Carat:       req.Carat,
		Color:       req.Color,
		Clarity:     req.Clarity,
		Origin:      req.Origin,
		Certificate: req.Certificate,
		ImageURL:    req.ImageURL,
		Status:      domain.GemAvailable,
	}

	if err := s.gemRepo.Create(g); err != nil {
		return nil, err
	}

	return g, nil
}

func (s *GemService) GetByID(id int64) (*domain.Gem, error) {
	if id <= 0 {
		return nil, errors.New("invalid gem id")
	}
	return s.gemRepo.GetByID(id)
}
