package player

// import (
// 	"fmt"

// 	"github.com/Zigelzi/go-tiimit/internal/db"
// )

// func GetAll() (players []Player, err error) {
// 	rows, err := db.DB.Query("SELECT id, name, myclub_id, run_power, ball_handling, is_goalie FROM players")
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get all players with error: %w", err)
// 	}
// 	defer rows.Close()

// 	rowNumber := 0
// 	for rows.Next() {
// 		var player Player
// 		rowNumber++
// 		err := rows.Scan(&player.id, &player.Name, &player.MyClubId, &player.runPower, &player.ballHandling, &player.IsGoalie)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to scan player on row %d: %w", rowNumber, err)
// 		}
// 		players = append(players, player)
// 	}
// 	return players, nil
// }

// func Get(myClubId int64) (player db.Player, err error) {
// 	err = db.DB.QueryRow("SELECT id, name, myclub_id, run_power, ball_handling, is_goalie FROM players WHERE myclub_id=?", myClubId).Scan(&player.Id, &player.Name, &player.MyClubId, &player.RunPower, &player.BallHandling, &player.IsGoalie)

// 	if err != nil {
// 		return db.Player{}, fmt.Errorf("unable to query player with MyClub ID %d: %w", myClubId, err)
// 	}
// 	return player, err
// }

// func Insert(player Player) error {
// 	_, err := db.DB.Exec("INSERT INTO players (name, myclub_id, run_power, ball_handling) VALUES (?, ?, ?, ?)", player.Name, player.MyClubId, player.runPower, player.ballHandling)
// 	return err
// }

// func ToggleGoalieStatus(player Player) error {
// 	query := `
// 	UPDATE players
// 	SET is_goalie = ?
// 	WHERE id=?
// 	`
// 	_, err := db.DB.Exec(query, !player.IsGoalie, player.id)
// 	if err != nil {
// 		return fmt.Errorf("unable to update goalies status of player %d: %w", player.id, err)
// 	}
// 	return nil
// }

// func exists(myClubId int64) (isExisting bool, err error) {
// 	query := "SELECT EXISTS (SELECT 1 FROM players WHERE myclub_id=?)"
// 	err = db.DB.QueryRow(query, myClubId).Scan(&isExisting)
// 	if err != nil {
// 		return isExisting, err
// 	}
// 	return isExisting, nil
// }

// func updateRunPower(myclubId int64, newRunPower float64) error {
// 	query := `
// 	UPDATE players
// 	SET run_power=?
// 	WHERE myclub_id=?
// 	`
// 	_, err := db.DB.Exec(query, newRunPower, myclubId)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
