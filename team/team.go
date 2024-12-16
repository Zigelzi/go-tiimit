package team

import (
	"example.com/go-tiimit/player"
)

type Team struct {
	name    string
	players []player.Player
}

func New(name string) Team {
	return Team{
		name: name,
	}
}
