package routes

import (
	"chartbattle-backend/internal/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	r.GET("/game/segment", handlers.GetSegment)
	r.POST("/game/submit", handlers.SubmitPrediction)
	r.GET("/leaderboard", handlers.GetLeaderboard)
	r.GET("/rank/:user_id", handlers.GetRank)
	r.GET("/stats/:user_id", handlers.GetStats)

}
