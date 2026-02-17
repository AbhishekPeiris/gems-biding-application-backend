package service

import (
	"errors"

	"github.com/boswin/gems-auction-backend/internal/domain"
	"github.com/boswin/gems-auction-backend/internal/repository"
)

type ChatService struct {
	chatRepo   *repository.ChatRepository
	broadcast  AuctionEventBroadcaster // reuse same broadcaster interface (can be nil)
}

func NewChatService(chatRepo *repository.ChatRepository, broadcast AuctionEventBroadcaster) *ChatService {
	return &ChatService{chatRepo: chatRepo, broadcast: broadcast}
}

type SendChatRequest struct {
	AuctionID int64  `json:"auction_id"`
	UserID    int64  `json:"user_id"`
	Message   string `json:"message"`
}

func (s *ChatService) Send(req SendChatRequest) (*domain.ChatMessage, error) {
	if req.AuctionID <= 0 || req.UserID <= 0 {
		return nil, errors.New("auction_id and user_id required")
	}
	if req.Message == "" {
		return nil, errors.New("message required")
	}

	msg := &domain.ChatMessage{
		AuctionID: req.AuctionID,
		UserID:    req.UserID,
		Message:   req.Message,
	}

	if err := s.chatRepo.Create(msg); err != nil {
		return nil, err
	}

	if s.broadcast != nil {
		s.broadcast.BroadcastToAuction(req.AuctionID, "CHAT_MESSAGE", msg)
	}

	return msg, nil
}

func (s *ChatService) GetByAuction(auctionID int64) ([]domain.ChatMessage, error) {
	if auctionID <= 0 {
		return nil, errors.New("invalid auction id")
	}
	return s.chatRepo.GetByAuction(auctionID)
}
