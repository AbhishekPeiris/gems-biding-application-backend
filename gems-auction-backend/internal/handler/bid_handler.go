package handler

import (
	"net/http"

	"github.com/boswin/gems-auction-backend/internal/service"
	"github.com/gin-gonic/gin"
)

type BidHandler struct {
	bidService *service.BidService
}

func NewBidHandler(bidService *service.BidService) *BidHandler {
	return &BidHandler{bidService: bidService}
}

func (h *BidHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("", h.PlaceBid)
}

func (h *BidHandler) PlaceBid(c *gin.Context) {
	var req service.PlaceBidRequest
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

	bid, err := h.bidService.PlaceBid(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, bid)
}
