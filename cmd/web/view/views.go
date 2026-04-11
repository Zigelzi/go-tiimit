package view

import (
	"fmt"
	"time"

	"github.com/Zigelzi/go-tiimit/internal/practice"
)

type Practice struct {
	Date  time.Time
	Teams []Team
}

type Team struct {
	Number     int
	TotalScore float64
	Players    []Player
}

type Player struct {
	ID            int64
	Name          string
	Score         float64
	IsGoalie      bool
	HasVest       bool
	MoveURL       string
	ToggleVestURL string
}

func FromPractice(players []practice.PracticePlayer, teamNumber int) Team {
	newTeam := Team{
		Number:     teamNumber,
		TotalScore: practice.TotalScore(players),
	}
	for _, player := range players {
		newTeam.Players = append(newTeam.Players, FromPlayer(player))
	}
	return newTeam
}

func (t *Team) GeneratePlayerURLs(practiceId int64) {
	for i := range t.Players {
		t.Players[i].GenerateURLs(practiceId)
	}
}

func (t *Team) VestCount() int {
	numberOfVests := 0

	for _, player := range t.Players {
		if player.HasVest {
			numberOfVests++
		}
	}
	return numberOfVests
}

func FromPlayer(p practice.PracticePlayer) Player {
	return Player{
		ID:       p.Player.ID,
		Name:     p.Player.Name,
		IsGoalie: p.Player.IsGoalie,
		HasVest:  p.HasVest,
		Score:    p.Player.Score(),
	}
}

func (p *Player) GenerateURLs(practiceId int64) {
	p.MoveURL = fmt.Sprintf("/practices/%d/players/%d", practiceId, p.ID)
	p.ToggleVestURL = fmt.Sprintf("/practices/%d/players/%d/vest", practiceId, p.ID)
}
