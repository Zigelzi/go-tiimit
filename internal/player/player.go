package player

import (
	"fmt"

	"github.com/Zigelzi/go-tiimit/internal/db"
)

const minScore = 0
const maxScore = 10

var ErrInvalidScore = fmt.Errorf("score must be between %d-%d", minScore, maxScore)

type Player struct {
	ID           int64
	MyClubId     int64
	Name         string
	runPower     float64
	ballHandling float64
	IsGoalie     bool
}

func New(myclub_id int64, name string, runPower float64, ballHandling float64, isGoalie bool) Player {
	if myclub_id < 0 {
		myclub_id = 0
	}

	runPower = clampScore(runPower)
	ballHandling = clampScore(ballHandling)

	return Player{
		MyClubId:     myclub_id,
		Name:         name,
		runPower:     runPower,
		ballHandling: ballHandling,
		IsGoalie:     isGoalie,
	}
}

func clampScore(score float64) float64 {
	if score < minScore {
		return minScore
	}
	if score > maxScore {
		return maxScore
	}

	return score
}

func ValidateScore(score float64) error {
	if score < minScore || score > maxScore {
		return ErrInvalidScore
	}

	return nil
}

func FromDB(dbPlayer db.Player) Player {
	newPlayer := New(
		dbPlayer.MyclubID,
		dbPlayer.Name,
		dbPlayer.RunPower,
		dbPlayer.BallHandling,
		dbPlayer.IsGoalie)

	newPlayer.ID = dbPlayer.ID
	return newPlayer
}

func (p Player) RunPower() float64 {
	return p.runPower
}
func (p Player) BallHandling() float64 {
	return p.ballHandling
}
