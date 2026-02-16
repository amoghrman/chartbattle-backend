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
		SegmentID  string `json:"segment_id"`
		Prediction string `json:"prediction"` // "UP" or "DOWN"
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

	if len(candles) == 0 || len(future) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid segment data",
		})
		return
	}

	lastClose := candles[len(candles)-1].Close
	futureClose := future[0].Close

	actualDirection := "DOWN"
	if futureClose > lastClose {
		actualDirection = "UP"
	}

	result := "wrong"
	if req.Prediction == actualDirection {
		result = "correct"
	}

	c.JSON(http.StatusOK, gin.H{
		"result":           result,
		"actual_direction": actualDirection,
	})
}
