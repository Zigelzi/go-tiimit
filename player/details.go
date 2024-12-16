package player

import "fmt"

func (player Player) Details() {
	fmt.Printf("[%d] %s score: %.1f\n", player.MyClubId, player.Name, player.Score())
}
