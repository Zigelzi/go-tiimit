package player

type Player struct {
	MyClubId     int64
	Name         string
	runPower     float64
	ballHandling float64
}

type ByScore []Player

func (p ByScore) Len() int           { return len(p) }
func (p ByScore) Less(i, j int) bool { return p[i].Score() > p[j].Score() }
func (p ByScore) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func New(id int64, name string, runPower float64, ballHandling float64) Player {
	return Player{
		MyClubId:     id,
		Name:         name,
		runPower:     runPower,
		ballHandling: ballHandling,
	}
}
