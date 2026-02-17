package handler

import (
	"net/http"
	"strconv"

	"github.com/boswin/gems-auction-backend/internal/domain"
	"github.com/boswin/gems-auction-backend/internal/service"
	"github.com/gin-gonic/gin"
)

type GemHandler struct {
	gemService *service.GemService
}

func NewGemHandler(gemService *service.GemService) *GemHandler {
	return &GemHandler{gemService: gemService}
}

func (h *GemHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("", h.CreateGem)
	rg.GET("/:id", h.GetGemByID)
}

func (h *GemHandler) CreateGem(c *gin.Context) {
	var req service.CreateGemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Prefer seller_id from token middleware if available
	if v, ok := c.Get("user_id"); ok && req.SellerID == 0 {
		if id, ok2 := v.(int64); ok2 {
			req.SellerID = id
		}
	}

	// (Optional) role check if middleware sets role
	if v, ok := c.Get("role"); ok {
		if roleStr, ok2 := v.(string); ok2 {
			if domain.UserRole(roleStr) != domain.RoleSeller && domain.UserRole(roleStr) != domain.RoleAdmin {
				c.JSON(http.StatusForbidden, gin.H{"error": "only SELLER/ADMIN can create gems"})
				return
			}
		}
	}

	gem, err := h.gemService.Create(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gem)
}

func (h *GemHandler) GetGemByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	gem, err := h.gemService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "gem not found"})
		return
	}

	c.JSON(http.StatusOK, gem)
}
