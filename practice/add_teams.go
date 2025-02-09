package practice

import (
	"github.com/Zigelzi/go-tiimit/team"
)

func (p *Practice) AddTeams(team1, team2 team.Team) error {

	p.Teams[0] = team1
	p.Teams[1] = team2
	return nil
}
