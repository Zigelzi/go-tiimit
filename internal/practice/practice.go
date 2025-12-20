package practice

import (
	"time"

	"github.com/Zigelzi/go-tiimit/internal/player"
)

type Practice struct {
	TeamOnePlayers []player.Player
	TeamTwoPlayers []player.Player
	UnknownPlayers []player.Player
	Date           time.Time
}
