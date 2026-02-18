package handler

import (
	"net/http"
	"strconv"

	"github.com/boswin/gems-auction-backend/internal/domain"
	"github.com/boswin/gems-auction-backend/internal/service"
	"github.com/gin-gonic/gin"
)

type AuctionHandler struct {
	auctionService *service.AuctionService
}

func NewAuctionHandler(auctionService *service.AuctionService) *AuctionHandler {
	return &AuctionHandler{auctionService: auctionService}
}

func (h *AuctionHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("", h.CreateAuction)
	rg.GET("/:id", h.GetAuctionByID)
	rg.POST("/:id/start", h.StartAuction)
	rg.POST("/:id/end", h.EndAuction)
}

func (h *AuctionHandler) CreateAuction(c *gin.Context) {
	// Optional role check
	if v, ok := c.Get("role"); ok {
		if roleStr, ok2 := v.(string); ok2 {
			role := domain.UserRole(roleStr)
			if role != domain.RoleSeller && role != domain.RoleAdmin {
				c.JSON(http.StatusForbidden, gin.H{"error": "only SELLER/ADMIN can create auctions"})
				return
			}
		}
	}

	var req service.CreateAuctionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	a, err := h.auctionService.Create(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, a)
}

func (h *AuctionHandler) GetAuctionByID(c *gin.Context) {
	auctionID, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	a, err := h.auctionService.GetByID(auctionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "auction not found"})
		return
	}

	c.JSON(http.StatusOK, a)
}

func (h *AuctionHandler) StartAuction(c *gin.Context) {
	// Optional role check
	if v, ok := c.Get("role"); ok {
		if roleStr, ok2 := v.(string); ok2 {
			role := domain.UserRole(roleStr)
			if role != domain.RoleSeller && role != domain.RoleAdmin {
				c.JSON(http.StatusForbidden, gin.H{"error": "only SELLER/ADMIN can start auctions"})
				return
			}
		}
	}

	auctionID, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	if err := h.auctionService.StartAuction(auctionID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "auction started"})
}

func (h *AuctionHandler) EndAuction(c *gin.Context) {
	// Optional role check
	if v, ok := c.Get("role"); ok {
		if roleStr, ok2 := v.(string); ok2 {
			role := domain.UserRole(roleStr)
			if role != domain.RoleSeller && role != domain.RoleAdmin {
				c.JSON(http.StatusForbidden, gin.H{"error": "only SELLER/ADMIN can end auctions"})
				return
			}
		}
	}

	auctionID, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	// winner_id is optional
	var body struct {
		WinnerID *int64 `json:"winner_id"`
	}
	_ = c.ShouldBindJSON(&body)

	if err := h.auctionService.EndAuction(auctionID, body.WinnerID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "auction ended", "winner_id": body.WinnerID})
}

func parseIDParam(c *gin.Context, param string) (int64, bool) {
	idStr := c.Param(param)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return 0, false
	}
	return id, true
}

func (h *AuctionHandler) GetAllAuctions(c *gin.Context) {
	auctions, err := h.auctionService.GetAllAuctions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, auctions)
}

