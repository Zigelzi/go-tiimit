package team

import "github.com/Zigelzi/go-tiimit/internal/player"

type Team struct {
	Name    string
	Players []player.Player
}

func New(name string) Team {
	return Team{
		Name: name,
	}
}
