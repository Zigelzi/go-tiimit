package player

import "sort"

type ByScore []Player

func (p ByScore) Len() int           { return len(p) }
func (p ByScore) Less(i, j int) bool { return p[i].Score() > p[j].Score() }
func (p ByScore) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func SortByScore(players []Player) {
	sort.Sort(ByScore(players))
}
