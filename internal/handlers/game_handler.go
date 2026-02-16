package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"chartbattle-backend/internal/services"
)

func GetSegment(c *gin.Context) {

	segment, err := services.GetRandomSegment()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch segment",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":      segment.ID,
		"candles": segment.Candles,
	})
}

func SubmitPrediction(c *gin.Context) {

	type Request struct {
		UserID     string `json:"user_id"`
		SegmentID  string `json:"segment_id"`
		Prediction string `json:"prediction"`
	}

	var req Request
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	segment, err := services.GetRandomSegment() // You can later create GetSegmentByID
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Segment fetch failed"})
		return
	}

	result, xp := services.EvaluatePrediction(segment, req.Prediction)

	if err := services.SaveGameSession(req.UserID, req.SegmentID, req.Prediction, result, xp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	bonusXP, err := services.UpdateUserXP(req.UserID, xp, result)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result":    result,
		"xp_earned": xp,
		"bonus_xp":  bonusXP,
	})
}
