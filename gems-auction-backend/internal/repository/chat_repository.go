package repository

import (
	"context"
	"time"

	"github.com/boswin/gems-auction-backend/config"
	"github.com/boswin/gems-auction-backend/internal/domain"
)

type ChatRepository struct{}

func NewChatRepository() *ChatRepository {
	return &ChatRepository{}
}

func (r *ChatRepository) Create(msg *domain.ChatMessage) error {
	query := `
		INSERT INTO chat_messages (auction_id,user_id,message,created_at)
		VALUES ($1,$2,$3,$4)
		RETURNING id
	`

	return config.DB.QueryRow(context.Background(), query,
		msg.AuctionID,
		msg.UserID,
		msg.Message,
		time.Now(),
	).Scan(&msg.ID)
}

func (r *ChatRepository) GetByAuction(auctionID int64) ([]domain.ChatMessage, error) {
	query := `
		SELECT id,auction_id,user_id,message,created_at
		FROM chat_messages
		WHERE auction_id=$1
		ORDER BY created_at ASC
	`

	rows, err := config.DB.Query(context.Background(), query, auctionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []domain.ChatMessage

	for rows.Next() {
		var msg domain.ChatMessage
		err := rows.Scan(
			&msg.ID,
			&msg.AuctionID,
			&msg.UserID,
			&msg.Message,
			&msg.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, nil
}
