package player

import (
	"fmt"

	"github.com/Zigelzi/go-tiimit/db"
)

func GetAll() (players []Player, err error) {
	rows, err := db.DB.Query("SELECT id, name, myclub_id, run_power, ball_handling FROM players")
	if err != nil {
		return nil, fmt.Errorf("failed to get all players with error: %w", err)
	}
	defer rows.Close()

	rowNumber := 0
	for rows.Next() {
		var player Player
		rowNumber++
		err := rows.Scan(&player.id, &player.Name, &player.MyClubId, &player.runPower, &player.ballHandling)
		if err != nil {
			return nil, fmt.Errorf("failed to scan player on row: %d", rowNumber)
		}
		players = append(players, player)
	}
	return players, nil
}

func Insert(player Player) error {
	_, err := db.DB.Exec("INSERT INTO players (name, myclub_id, run_power, ball_handling) VALUES (?, ?, ?, ?)", player.Name, player.MyClubId, player.runPower, player.ballHandling)
	return err
}
