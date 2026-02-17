package domain

import "time"

type GemStatus string

const (
	GemAvailable GemStatus = "AVAILABLE"
	GemAuction   GemStatus = "AUCTION"
	GemSold      GemStatus = "SOLD"
)

type Gem struct {
	ID           int64     `json:"id"`
	SellerID     int64     `json:"seller_id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Carat        float64   `json:"carat"`
	Color        string    `json:"color"`
	Clarity      string    `json:"clarity"`
	Origin       string    `json:"origin"`
	Certificate  string    `json:"certificate"` // certificate number
	ImageURL     string    `json:"image_url"`
	Status       GemStatus `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
