package services

import (
	"encoding/json"
	// "errors"

	"chartbattle-backend/internal/database"
	"chartbattle-backend/internal/models"
)

func GetRandomSegment() (*models.ChartSegment, error) {
	row := database.DB.QueryRow(`
		SELECT id, candles, future
		FROM chart_segments
		ORDER BY RANDOM()
		LIMIT 1
	`)

	var segment models.ChartSegment
	var candlesStr string
	var futureStr string

	err := row.Scan(&segment.ID, &candlesStr, &futureStr)
	if err != nil {
		return nil, err
	}

	json.Unmarshal([]byte(candlesStr), &segment.Candles)
	json.Unmarshal([]byte(futureStr), &segment.Future)

	return &segment, nil
}

func EvaluatePrediction(segment *models.ChartSegment, prediction string) (string, int) {

	if len(segment.Candles) == 0 || len(segment.Future) == 0 {
		return "invalid", 0
	}

	lastClose := segment.Candles[len(segment.Candles)-1].Close
	futureClose := segment.Future[0].Close

	actualDirection := "DOWN"
	if futureClose > lastClose {
		actualDirection = "UP"
	}

	if prediction == actualDirection {
		return "correct", 100
	}

	return "wrong", 0
}

func SaveGameSession(userID, segmentID, prediction, result string, xp int) error {

	_, err := database.DB.Exec(`
		INSERT INTO game_sessions (user_id, segment_id, prediction, result, xp_earned)
		VALUES ($1, $2, $3, $4, $5)
	`, userID, segmentID, prediction, result, xp)

	return err
}
