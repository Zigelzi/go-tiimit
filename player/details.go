package player

import "fmt"

func (player Player) Details() (details string) {
	if player.IsGoalie {
		goalieSymbol := "[G]"
		return fmt.Sprintf("%s %s", player.Name, goalieSymbol)
	} else {
		return player.Name
	}
}
