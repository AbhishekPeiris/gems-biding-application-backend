package main

import (
	"log"
	"net/http"
	"time"

	"github.com/boswin/gems-auction-backend/config"
	"github.com/boswin/gems-auction-backend/internal/handler"
	"github.com/boswin/gems-auction-backend/internal/middleware"
	"github.com/boswin/gems-auction-backend/internal/repository"
	"github.com/boswin/gems-auction-backend/internal/service"
	"github.com/boswin/gems-auction-backend/internal/websocket"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	// ===============================
	// 1️⃣ Load Configuration
	// ===============================
	config.LoadConfig()

	// ===============================
	// 2️⃣ Connect Database
	// ===============================
	config.ConnectDatabase()
	defer config.CloseDatabase()

	// ===============================
	// 3️⃣ Initialize WebSocket Manager
	// ===============================
	wsManager := websocket.NewManager()

	// ===============================
	// 4️⃣ Initialize Repositories
	// ===============================
	userRepo := repository.NewUserRepository()
	gemRepo := repository.NewGemRepository()
	auctionRepo := repository.NewAuctionRepository()
	bidRepo := repository.NewBidRepository()
	chatRepo := repository.NewChatRepository()

	// ===============================
	// 5️⃣ Initialize Services
	// ===============================
	authService := service.NewAuthService(userRepo)
	gemService := service.NewGemService(gemRepo)
	auctionService := service.NewAuctionService(auctionRepo)
	bidService := service.NewBidService(bidRepo, auctionRepo, wsManager)
	chatService := service.NewChatService(chatRepo, wsManager)
	paymentService := service.NewPaymentService()

	_ = paymentService

	// ===============================
	// 6️⃣ Initialize Handlers
	// ===============================
	authHandler := handler.NewAuthHandler(authService)
	gemHandler := handler.NewGemHandler(gemService)
	auctionHandler := handler.NewAuctionHandler(auctionService)
	bidHandler := handler.NewBidHandler(bidService)
	chatHandler := handler.NewChatHandler(chatService)
	wsHandler := handler.NewWebSocketHandler(wsManager)

	// ===============================
	// 7️⃣ Setup Gin Router
	// ===============================
	r := gin.New()

	// Logging + Recovery
	r.Use(middleware.LoggingMiddleware())
	r.Use(gin.Recovery())

	// ===============================
	// 🔥 CORS CONFIGURATION (IMPORTANT FIX)
	// ===============================
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// ===============================
	// 8️⃣ Health Check
	// ===============================
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "running",
			"service": "gems-auction-backend",
			"time":    time.Now(),
		})
	})

	// ===============================
	// 9️⃣ WebSocket Route
	// ===============================
	wsHandler.RegisterRoutes(r)

	// ===============================
	// 🔟 API ROUTES
	// ===============================
	api := r.Group("/api")

	// -------- AUTH (Public) --------
	authGroup := api.Group("/auth")
	authHandler.RegisterRoutes(authGroup)

	// -------- Protected Routes --------
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware())

	// =====================================
	// GEMS ROUTES
	// =====================================
	gems := protected.Group("/gems")

	gems.POST("",
		middleware.RoleMiddleware("SELLER", "ADMIN"),
		gemHandler.CreateGem,
	)

	gems.GET("/:id", gemHandler.GetGemByID)

	// =====================================
	// AUCTION ROUTES
	// =====================================
	auctions := protected.Group("/auctions")

	auctions.POST("",
		middleware.RoleMiddleware("SELLER", "ADMIN"),
		auctionHandler.CreateAuction,
	)

	auctions.GET("/:id", auctionHandler.GetAuctionByID)

	auctions.POST("/:id/start",
		middleware.RoleMiddleware("SELLER", "ADMIN"),
		auctionHandler.StartAuction,
	)

	auctions.POST("/:id/end",
		middleware.RoleMiddleware("SELLER", "ADMIN"),
		auctionHandler.EndAuction,
	)

	// =====================================
	// BIDDING ROUTES
	// =====================================
	bids := protected.Group("/bids")

	bids.POST("",
		middleware.RoleMiddleware("BUYER", "ADMIN"),
		bidHandler.PlaceBid,
	)

	// =====================================
	// CHAT ROUTES
	// =====================================
	chat := protected.Group("/chat")

	chat.POST("", chatHandler.SendChat)
	chat.GET("/auction/:id", chatHandler.GetChatByAuction)

	// ===============================
	// 🚀 Start Server
	// ===============================
	port := config.AppConfig.Port
	log.Println("🚀 Gems Auction Backend running on port:", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
