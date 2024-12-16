package practice

import "example.com/go-tiimit/team"

func (p *Practice) CreateTeams() error {
	team1, team2, err := team.Distribute(p.Players)
	if err != nil {
		return err
	}
	p.Teams[0] = team1
	p.Teams[1] = team2
	return nil
}
