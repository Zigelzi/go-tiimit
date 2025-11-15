package db

import (
	"fmt"
)

func (q *Queries) Get(myClubId int64) (player Player, err error) {
	err = q.db.QueryRow("SELECT id, name, myclub_id, run_power, ball_handling, is_goalie FROM players WHERE myclub_id=?", myClubId).Scan(&player.Id, &player.Name, &player.MyClubId, &player.RunPower, &player.BallHandling, &player.IsGoalie)

	if err != nil {
		return Player{}, fmt.Errorf("unable to query player with MyClub ID %d: %w", myClubId, err)
	}
	return player, err
}
