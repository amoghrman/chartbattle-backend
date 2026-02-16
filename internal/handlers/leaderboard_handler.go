package handlers

import (
	"net/http"

	"chartbattle-backend/internal/database"
	"github.com/gin-gonic/gin"
)

func GetLeaderboard(c *gin.Context) {

	rows, err := database.DB.Query(`
		SELECT email, xp, total_games
		FROM users
		ORDER BY xp DESC
		LIMIT 10
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch leaderboard",
		})
		return
	}
	defer rows.Close()

	type LeaderboardUser struct {
		Email      string `json:"email"`
		XP         int    `json:"xp"`
		TotalGames int    `json:"total_games"`
	}

	var leaderboard []LeaderboardUser

	for rows.Next() {
		var user LeaderboardUser
		rows.Scan(&user.Email, &user.XP, &user.TotalGames)
		leaderboard = append(leaderboard, user)
	}

	c.JSON(http.StatusOK, gin.H{
		"leaderboard": leaderboard,
	})
}
