package player

type Player struct {
	id           int64
	MyClubId     int64
	Name         string
	runPower     float64
	ballHandling float64
	IsGoalie     bool
}

type ByScore []Player

func (p ByScore) Len() int           { return len(p) }
func (p ByScore) Less(i, j int) bool { return p[i].Score() > p[j].Score() }
func (p ByScore) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func New(myclub_id int64, name string, runPower float64, ballHandling float64, isGoalie bool) Player {
	return Player{
		MyClubId:     myclub_id,
		Name:         name,
		runPower:     runPower,
		ballHandling: ballHandling,
		IsGoalie:     isGoalie,
	}
}
