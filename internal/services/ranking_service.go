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
