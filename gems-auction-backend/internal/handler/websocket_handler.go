package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// WSManager is implemented by your websocket hub layer.
// We'll implement this in internal/websocket later.
type WSManager interface {
	ServeAuctionWS(c *gin.Context, auctionID int64)
}

type WebSocketHandler struct {
	ws WSManager
}

func NewWebSocketHandler(ws WSManager) *WebSocketHandler {
	return &WebSocketHandler{ws: ws}
}

func (h *WebSocketHandler) RegisterRoutes(r *gin.Engine) {
	// WebSocket endpoint (not inside /api usually)
	r.GET("/ws/auction/:id", h.HandleAuctionWS)
}

func (h *WebSocketHandler) HandleAuctionWS(c *gin.Context) {
	idStr := c.Param("id")
	auctionID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || auctionID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid auction id"})
		return
	}

	h.ws.ServeAuctionWS(c, auctionID)
}
