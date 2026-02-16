package services

import (
	"chartbattle-backend/internal/database"
)

func UpdateUserXP(userID string, baseXP int, result string) (int, error) {

	var currentStreak int
	var bestStreak int
	var currentRank string

	row := database.DB.QueryRow(`
		SELECT streak, best_streak, rank
		FROM users
		WHERE id = $1
	`, userID)

	err := row.Scan(&currentStreak, &bestStreak, &currentRank)
	if err != nil {
		return 0, err
	}

	bonusXP := 0
	totalXPChange := baseXP

	if result == "correct" {
		currentStreak++

		// streak bonus
		if currentStreak == 3 {
			bonusXP = 20
		} else if currentStreak == 5 {
			bonusXP = 50
		}

		if currentStreak > bestStreak {
			bestStreak = currentStreak
		}

		totalXPChange = baseXP + bonusXP

	} else {
		currentStreak = 0

		// rank-based loss
		switch currentRank {
		case "Trader":
			totalXPChange = -20
		case "Sniper":
			totalXPChange = -40
		case "Alpha":
			totalXPChange = -70
		default:
			totalXPChange = 0
		}
	}

	// Update user core fields
	_, err = database.DB.Exec(`
		UPDATE users
		SET xp = GREATEST(xp + $1, 0),
		    total_games = total_games + 1,
		    streak = $2,
		    best_streak = $3,
		    total_correct = total_correct + CASE WHEN $4 = 'correct' THEN 1 ELSE 0 END,
		    total_wrong = total_wrong + CASE WHEN $4 = 'wrong' THEN 1 ELSE 0 END
		WHERE id = $5
	`, totalXPChange, currentStreak, bestStreak, result, userID)

	if err != nil {
		return 0, err
	}

	// Recalculate rank
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
