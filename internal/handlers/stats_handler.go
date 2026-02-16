package handlers

import (
	"net/http"

	"chartbattle-backend/internal/services"
	"github.com/gin-gonic/gin"
)

func GetStats(c *gin.Context) {

	userID := c.Param("user_id")

	stats, err := services.GetUserStats(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"stats": stats,
	})
}
