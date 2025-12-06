package practice

import (
	"github.com/Zigelzi/go-tiimit/internal/team"
)

type Practice struct {
	AttendingPlayers map[int]AttendanceStatus
	Teams            [2]team.Team
}

func New() Practice {
	return Practice{
		AttendingPlayers: make(map[int]AttendanceStatus),
	}
}
