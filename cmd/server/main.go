package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"chartbattle-backend/internal/config"
	"chartbattle-backend/internal/database"
	"chartbattle-backend/internal/routes"
)

func main() {
	// Load environment variables
	config.LoadConfig()

	// Connect database
	database.Connect()

	// Initialize Gin
	router := gin.Default()

	// Avoid proxy warning (optional)
	router.SetTrustedProxies(nil)

	// Register routes
	routes.RegisterRoutes(router)

	// Start server
	port := config.GetEnv("PORT", "8080")
	log.Println("Server running on port", port)
	router.Run(":" + port)
}
