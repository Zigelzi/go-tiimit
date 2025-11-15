package db

type Player struct {
	Id           int64
	MyClubId     int64
	Name         string
	RunPower     float64
	BallHandling float64
	IsGoalie     bool
}
