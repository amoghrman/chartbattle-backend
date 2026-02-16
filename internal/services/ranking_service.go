package services

import "chartbattle-backend/internal/database"

type LeaderboardUser struct {
	Email      string `json:"email"`
	XP         int    `json:"xp"`
	Rank       string `json:"rank"`
	TotalGames int    `json:"total_games"`
}

func GetLeaderboard() ([]LeaderboardUser, error) {

	rows, err := database.DB.Query(`
		SELECT email, xp, rank, total_games
		FROM users
		ORDER BY xp DESC
		LIMIT 10
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var leaderboard []LeaderboardUser

	for rows.Next() {
		var user LeaderboardUser
		rows.Scan(&user.Email, &user.XP, &user.Rank, &user.TotalGames)
		leaderboard = append(leaderboard, user)
	}

	return leaderboard, nil
}

type UserStats struct {
	Email         string  `json:"email"`
	XP            int     `json:"xp"`
	Rank          string  `json:"rank"`
	TotalGames    int     `json:"total_games"`
	TotalCorrect  int     `json:"total_correct"`
	TotalWrong    int     `json:"total_wrong"`
	Accuracy      float64 `json:"accuracy"`
	CurrentStreak int     `json:"current_streak"`
	BestStreak    int     `json:"best_streak"`
	NextRank      string  `json:"next_rank"`
	XPToNextRank  int     `json:"xp_to_next_rank"`
}

func GetUserStats(userID string) (UserStats, error) {

	row := database.DB.QueryRow(`
		SELECT email, xp, rank, total_games,
		       total_correct, total_wrong,
		       streak, best_streak
		FROM users
		WHERE id = $1
	`, userID)

	var stats UserStats
	var streak int

	err := row.Scan(
		&stats.Email,
		&stats.XP,
		&stats.Rank,
		&stats.TotalGames,
		&stats.TotalCorrect,
		&stats.TotalWrong,
		&streak,
		&stats.BestStreak,
	)
	if err != nil {
		return stats, err
	}

	stats.CurrentStreak = streak

	// Calculate accuracy
	if stats.TotalGames > 0 {
		stats.Accuracy = float64(stats.TotalCorrect) / float64(stats.TotalGames) * 100
	}

	// Next rank logic
	switch stats.Rank {
	case "Rookie":
		stats.NextRank = "Trader"
		stats.XPToNextRank = 500 - stats.XP
	case "Trader":
		stats.NextRank = "Sniper"
		stats.XPToNextRank = 1000 - stats.XP
	case "Sniper":
		stats.NextRank = "Alpha"
		stats.XPToNextRank = 2000 - stats.XP
	default:
		stats.NextRank = "Max"
		stats.XPToNextRank = 0
	}

	if stats.XPToNextRank < 0 {
		stats.XPToNextRank = 0
	}

	return stats, nil
}
func GetUserRank(userID string) (LeaderboardUser, error) {

	row := database.DB.QueryRow(`
		SELECT email, xp, rank, total_games
		FROM users
		WHERE id = $1
	`, userID)

	var user LeaderboardUser

	err := row.Scan(&user.Email, &user.XP, &user.Rank, &user.TotalGames)

	return user, err
}
