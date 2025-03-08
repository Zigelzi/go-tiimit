package practice

import (
	"github.com/Zigelzi/go-tiimit/player"
	"github.com/Zigelzi/go-tiimit/team"
)

type Practice struct {
	AttendingPlayers map[int]AttendanceStatus
	Players          []player.Player
	Teams            [2]team.Team
}

func New() Practice {
	return Practice{
		AttendingPlayers: make(map[int]AttendanceStatus),
	}
}
