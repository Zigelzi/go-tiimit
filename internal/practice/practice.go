package practice

import (
	"errors"
	"math"
	"sort"
	"time"

	"github.com/Zigelzi/go-tiimit/internal/db"
	"github.com/Zigelzi/go-tiimit/internal/player"
)

var ErrNoPracticeRows = errors.New("no practice rows")

type Practice struct {
	ID             int64
	TeamOnePlayers []PracticePlayer
	TeamTwoPlayers []PracticePlayer
	UnknownPlayers []PracticePlayer
	Date           time.Time
}

type PracticePlayer struct {
	Player  player.Player
	HasVest bool
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

	teamOnePlayers := []PracticePlayer{}
	teamTwoPlayers := []PracticePlayer{}
	for _, row := range dbPracticeRows {
		// TODO: Reuse creating the practice players to improve maintainability.
		currentPlayer := PracticePlayer{
			Player: player.New(
				row.MyclubID.Int64,
				row.Name.String,
				row.RunPower.Float64,
				row.BallHandling.Float64,
				row.IsGoalie.Bool,
			),
			HasVest: row.HasVest,
		}

		currentPlayer.Player.ID = row.PlayerID.Int64

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

func FromPlayer(players []player.Player) []PracticePlayer {
	practicePlayers := make([]PracticePlayer, len(players))
	for i, player := range players {
		practicePlayers[i] = PracticePlayer{
			Player:  player,
			HasVest: false,
		}
	}
	return practicePlayers
}

func PracticePlayerFromDB(dbPracticePlayer db.GetPracticeTeamPlayersRow) PracticePlayer {
	player := player.New(
		dbPracticePlayer.MyclubID,
		dbPracticePlayer.Name,
		dbPracticePlayer.RunPower,
		dbPracticePlayer.BallHandling,
		dbPracticePlayer.IsGoalie,
	)
	player.ID = dbPracticePlayer.PlayerID
	return PracticePlayer{
		Player:  player,
		HasVest: dbPracticePlayer.HasVest,
	}
}

func SortByScore(players []PracticePlayer) {
	sort.Slice(players, func(i, j int) bool {
		return players[i].Player.Score() > players[j].Player.Score()
	})
}

func TotalScore(players []PracticePlayer) float64 {
	totalScore := 0.0
	for _, player := range players {
		totalScore += player.Player.Score()
	}
	return roundToTwoDecimals(totalScore)
}

func roundToTwoDecimals(value float64) float64 {
	return math.Round(value*100) / 100
}
