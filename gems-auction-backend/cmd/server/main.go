package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/boswin/gems-auction-backend/config"
)

func main() {
	// Load configuration
	config.LoadConfig()

	// Connect to database
	config.ConnectDatabase()
	defer config.CloseDatabase()

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Gems Auction Backend Running 🚀",
		})
	})

	port := config.AppConfig.Port

	log.Println("🚀 Server running on port:", port)
	r.Run(":" + port)
}
