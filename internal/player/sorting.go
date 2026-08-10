package player

import "sort"

const (
	SortKeyRunPower     = "run-power"
	SortKeyBallHandling = "ball-handling"
	SortKeyScore        = "score"
)

func SortPlayersByDescending(players []Player, sortKey string) {
	switch sortKey {
	case SortKeyRunPower:
		SortByRunpower(players)
	case SortKeyBallHandling:
		SortByBallHandling(players)
	case SortKeyScore:
		SortByScore(players)
	}
}

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
