package models

import (
	"time"

	"github.com/google/uuid"
)

type GameSession struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	SegmentID  uuid.UUID `json:"segment_id"`
	Prediction string    `json:"prediction"`
	Result     string    `json:"result"`
	XPEarned   int       `json:"xp_earned"`
	CreatedAt  time.Time `json:"created_at"`
}
