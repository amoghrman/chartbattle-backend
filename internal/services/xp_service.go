package services

import (
	"chartbattle-backend/internal/database"
)

func UpdateUserXP(userID string, baseXP int, result string) (int, error) {

	var currentStreak int
	var currentRank string

	// Get current streak and rank
	row := database.DB.QueryRow(`
		SELECT streak, rank FROM users WHERE id = $1
	`, userID)

	err := row.Scan(&currentStreak, &currentRank)
	if err != nil {
		return 0, err
	}

	bonusXP := 0
	totalXPChange := baseXP

	if result == "correct" {
		currentStreak++

		// Streak bonus
		if currentStreak == 3 {
			bonusXP = 20
		} else if currentStreak == 5 {
			bonusXP = 50
		}

		totalXPChange = baseXP + bonusXP

	} else {
		// Reset streak
		currentStreak = 0

		// Rank-based XP loss
		switch currentRank {
		case "Trader":
			totalXPChange = -20
		case "Sniper":
			totalXPChange = -40
		case "Alpha":
			totalXPChange = -70
		default: // Rookie
			totalXPChange = 0
		}
	}

	// Update XP safely (never below 0)
	_, err = database.DB.Exec(`
		UPDATE users
		SET xp = GREATEST(xp + $1, 0),
		    total_games = total_games + 1,
		    streak = $2
		WHERE id = $3
	`, totalXPChange, currentStreak, userID)

	if err != nil {
		return 0, err
	}

	// Recalculate rank after XP update
	_, err = database.DB.Exec(`
		UPDATE users
		SET rank = CASE
			WHEN xp >= 2000 THEN 'Alpha'
			WHEN xp >= 1000 THEN 'Sniper'
			WHEN xp >= 500 THEN 'Trader'
			ELSE 'Rookie'
		END
		WHERE id = $1
	`, userID)

	return bonusXP, err
}
