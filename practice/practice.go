package practice

import (
	"example.com/go-tiimit/player"
	"example.com/go-tiimit/team"
)

type Practice struct {
	id      int64
	Players []player.Player
	Teams   [2]team.Team
}

func New() Practice {
	return Practice{
		id: 1,
	}
}
