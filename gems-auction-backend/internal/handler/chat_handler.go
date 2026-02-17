package handler

import (
	"net/http"

	"github.com/boswin/gems-auction-backend/internal/service"
	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	chatService *service.ChatService
}

func NewChatHandler(chatService *service.ChatService) *ChatHandler {
	return &ChatHandler{chatService: chatService}
}

func (h *ChatHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("", h.SendChat)
	// chat history by auction
	rg.GET("/auction/:id", h.GetChatByAuction)
}

func (h *ChatHandler) SendChat(c *gin.Context) {
	var req service.SendChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Prefer user_id from token middleware if available
	if v, ok := c.Get("user_id"); ok {
		if id, ok2 := v.(int64); ok2 {
			req.UserID = id
		}
	}

	msg, err := h.chatService.Send(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, msg)
}

func (h *ChatHandler) GetChatByAuction(c *gin.Context) {
	auctionID, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	msgs, err := h.chatService.GetByAuction(auctionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"auction_id": auctionID, "messages": msgs})
}
