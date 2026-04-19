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

type Player struct {
	ID            int64
	MyclubId      int64
	Name          string
	Score         float64
	RunPower      float64
	BallHandling  float64
	IsGoalie      bool
	HasVest       bool
	MoveURL       string
	ToggleVestURL string
}

func (p *Player) GenerateURLs(practiceId int64) {
	p.MoveURL = fmt.Sprintf("/practices/%d/players/%d", practiceId, p.ID)
	p.ToggleVestURL = fmt.Sprintf("/practices/%d/players/%d/vest", practiceId, p.ID)
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

func FromPlayer(p practice.PracticePlayer) Player {
	return Player{
		ID:       p.Player.ID,
		MyclubId: p.Player.MyClubId,
		Name:     p.Player.Name,
		IsGoalie: p.Player.IsGoalie,
		HasVest:  p.HasVest,
		Score:    p.Player.Score(),
	}
}

type AddPlayerForm struct {
	MyClubId     string
	Name         string
	FieldErrors  map[string]string
	GeneralError string
}

func (pf AddPlayerForm) HasError(field string) bool {
	_, ok := pf.FieldErrors[field]
	return ok
}
