package routes

import (
	"chartbattle-backend/internal/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// Game routes
	r.GET("/game/segment", handlers.GetSegment)
	r.POST("/game/submit", handlers.SubmitPrediction)
}
