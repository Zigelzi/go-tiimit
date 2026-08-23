package view

import (
	"fmt"
	"math"
	"time"

	"github.com/Zigelzi/go-tiimit/internal/player"
	"github.com/Zigelzi/go-tiimit/internal/practice"
)

type Practice struct {
	ID    int64
	Date  time.Time
	Teams []Team
}

func (p Practice) TotalPlayers() int {
	return len(p.Teams[0].Players) + len(p.Teams[1].Players)
}

type PracticeSummary struct {
	ID          int64
	Date        time.Time
	PlayerCount int64
}

func (ps PracticeSummary) DaysInPast() float64 {
	hoursInPast := time.Since(ps.Date).Hours()
	daysInPast := math.Floor(hoursInPast / 24)
	return daysInPast
}

func (ps PracticeSummary) IsUpcoming() bool {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, ps.Date.Location())
	return ps.Date.Compare(today) >= 0
}

type Team struct {
	Number     int
	TotalScore float64
	Players    []Player
}

func (t Team) VestBorderClass() string {
	if t.Number == 1 {
		return "border-vest-yellows"
	}
	return "border-vest-bibs"
}

func (t Team) VestBackgroundClass() string {
	if t.Number == 1 {
		return "bg-vest-yellows"
	}
	return "bg-vest-bibs"
}

func FromPractice(players []practice.PracticePlayer, teamNumber int) Team {
	newTeam := Team{
		Number:     teamNumber,
		TotalScore: practice.TotalScore(players),
	}
	for _, player := range players {
		newTeam.Players = append(newTeam.Players, FromPracticePlayer(player))
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

func (t *Team) TabURL(practiceId int64) string {
	return fmt.Sprintf("/practices/%d/teams/%d", practiceId, t.Number)
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

func FromPlayer(player player.Player) Player {
	return Player{
		ID:           player.ID,
		MyclubId:     player.MyClubId,
		Name:         player.Name,
		RunPower:     player.RunPower(),
		BallHandling: player.BallHandling(),
		IsGoalie:     player.IsGoalie,
		Score:        player.Score(),
	}
}

func (p *Player) ViewURL() string {
	return fmt.Sprintf("/players/%d", p.ID)
}

func (p *Player) EditURL() string {
	return fmt.Sprintf("/players/%d/edit", p.ID)
}

func (p *Player) SaveURL() string {
	// Return same value but to keep name explicit
	return p.ViewURL()
}

func (p *Player) RowID() string {
	return fmt.Sprintf("player-row-%d", p.ID)
}

func (p *Player) GenerateURLs(practiceId int64) {
	p.MoveURL = fmt.Sprintf("/practices/%d/players/%d", practiceId, p.ID)
	p.ToggleVestURL = fmt.Sprintf("/practices/%d/players/%d/vest", practiceId, p.ID)
}

func FromPracticePlayer(p practice.PracticePlayer) Player {
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

type EditPlayerRow struct {
	Player
	FieldErrors  map[string]string
	GeneralError string
}

func (epr EditPlayerRow) HasError(field string) bool {
	_, ok := epr.FieldErrors[field]
	return ok
}
