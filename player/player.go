package player

type Player struct {
	id           int64
	MyClubId     int64
	Name         string
	runPower     float64
	ballHandling float64
	IsGoalie     bool
}

func New(myclub_id int64, name string, runPower float64, ballHandling float64, isGoalie bool) Player {
	return Player{
		MyClubId:     myclub_id,
		Name:         name,
		runPower:     runPower,
		ballHandling: ballHandling,
		IsGoalie:     isGoalie,
	}
}
