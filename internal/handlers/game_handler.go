package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"

	"chartbattle-backend/internal/database"
)

type Candle struct {
	Open  float64 `json:"open"`
	Close float64 `json:"close"`
}

// GET /game/segment
func GetSegment(c *gin.Context) {
	row := database.DB.QueryRow(`
		SELECT id, candles
		FROM chart_segments
		ORDER BY RANDOM()
		LIMIT 1
	`)

	var id string
	var candlesStr string

	err := row.Scan(&id, &candlesStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch segment",
		})
		return
	}

	var candles []Candle
	json.Unmarshal([]byte(candlesStr), &candles)

	c.JSON(http.StatusOK, gin.H{
		"id":      id,
		"candles": candles,
	})
}

// POST /game/submit
func SubmitPrediction(c *gin.Context) {

	type Request struct {
		UserID     string `json:"user_id"`
		SegmentID  string `json:"segment_id"`
		Prediction string `json:"prediction"`
	}

	var req Request
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	row := database.DB.QueryRow(`
		SELECT candles, future
		FROM chart_segments
		WHERE id = $1
	`, req.SegmentID)

	var candlesStr string
	var futureStr string

	err := row.Scan(&candlesStr, &futureStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Segment not found",
		})
		return
	}

	var candles []Candle
	var future []Candle

	json.Unmarshal([]byte(candlesStr), &candles)
	json.Unmarshal([]byte(futureStr), &future)

	lastClose := candles[len(candles)-1].Close
	futureClose := future[0].Close

	actualDirection := "DOWN"
	if futureClose > lastClose {
		actualDirection = "UP"
	}

	result := "wrong"
	xpEarned := 0

	if req.Prediction == actualDirection {
		result = "correct"
		xpEarned = 100
	}

	// TEMPORARY: Hardcoded test user ID
	// testUserID := "a7f335ae-0f65-4382-b810-6c8faa88c95f" changed testuserID to taking input from the request body

	// Insert game session
	_, err = database.DB.Exec(`
		INSERT INTO game_sessions (user_id, segment_id, prediction, result, xp_earned)
		VALUES ($1, $2, $3, $4, $5)
	`,
		req.UserID,
		req.SegmentID,
		req.Prediction,
		result,
		xpEarned,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Update user XP
	_, err = database.DB.Exec(`
		UPDATE users
		SET xp = xp + $1,
			total_games = total_games + 1
		WHERE id = $2
	`,
		xpEarned,
		req.UserID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update user XP",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result":           result,
		"actual_direction": actualDirection,
		"xp_earned":        xpEarned,
	})
}
