package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID `json:"id"`
	Email      string    `json:"email"`
	XP         int       `json:"xp"`
	Rank       string    `json:"rank"`
	TotalGames int       `json:"total_games"`
	CreatedAt  time.Time `json:"created_at"`
}
