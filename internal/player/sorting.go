package player

import "sort"

func SortByScore(players []Player) {
	sort.Slice(players, func(i, j int) bool {
		return players[i].Score() > players[j].Score()
	})
}

func SortByRunpower(players []Player) {
	sort.Slice(players, func(i, j int) bool {
		return players[i].runPower > players[j].runPower
	})
}
func SortByBallHandling(players []Player) {
	sort.Slice(players, func(i, j int) bool {
		return players[i].ballHandling > players[j].ballHandling
	})
}
