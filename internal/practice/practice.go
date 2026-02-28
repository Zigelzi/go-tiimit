package practice

import (
	"errors"
	"time"

	"github.com/Zigelzi/go-tiimit/internal/db"
	"github.com/Zigelzi/go-tiimit/internal/player"
)

var ErrNoPracticeRows = errors.New("no practice rows")

type Practice struct {
	ID             int64
	TeamOnePlayers []player.Player
	TeamTwoPlayers []player.Player
	UnknownPlayers []player.Player
	Date           time.Time
}

func FromDB(dbPractice db.Practice) Practice {
	return Practice{
		ID:   dbPractice.ID,
		Date: dbPractice.Date,
	}
}

func FromDBWithPlayers(dbPracticeRows []db.GetPracticeWithPlayersRow) (Practice, error) {
	if len(dbPracticeRows) == 0 {
		return Practice{}, ErrNoPracticeRows
	}

	teamOnePlayers := []player.Player{}
	teamTwoPlayers := []player.Player{}
	for _, row := range dbPracticeRows {
		currentPlayer := player.New(
			row.MyclubID.Int64,
			row.Name.String,
			row.RunPower.Float64,
			row.BallHandling.Float64,
			row.IsGoalie.Bool,
		)
		if row.TeamNumber == 1 {
			teamOnePlayers = append(teamOnePlayers, currentPlayer)
		} else if row.TeamNumber == 2 {
			teamTwoPlayers = append(teamTwoPlayers, currentPlayer)
		}
	}
	newPractice := Practice{
		ID:             dbPracticeRows[0].PracticeID.Int64,
		Date:           dbPracticeRows[0].Date.Time,
		TeamOnePlayers: teamOnePlayers,
		TeamTwoPlayers: teamTwoPlayers,
	}
	return newPractice, nil
}
