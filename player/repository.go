package player

import "github.com/Zigelzi/go-tiimit/db"

func Insert(player Player) error {
	_, err := db.DB.Exec("INSERT INTO players (name, myclub_id, run_power, ball_handling) VALUES (?, ?, ?, ?)", player.Name, player.MyClubId, player.runPower, player.ballHandling)
	return err
}
