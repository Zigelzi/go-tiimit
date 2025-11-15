package practice

func (p *Practice) PrintTeams() {
	for _, team := range p.Teams {
		team.Details()
	}
}
