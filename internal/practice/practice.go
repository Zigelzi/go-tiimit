package practice

import (
	"time"

	"github.com/Zigelzi/go-tiimit/internal/player"
)

type Practice struct {
	ID             int64
	TeamOnePlayers []player.Player
	TeamTwoPlayers []player.Player
	UnknownPlayers []player.Player
	Date           time.Time
}
