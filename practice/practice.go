package practice

import (
	"github.com/Zigelzi/go-tiimit/player"
	"github.com/Zigelzi/go-tiimit/team"
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
