package file

import (
	"fmt"
	"strconv"
)

type ClubPlayerRow struct {
	PlayerRow    PlayerRow
	BallHandling float64
	RunPower     float64
}

func newClubPlayerRow(myClubId, name, runPower, ballHandling string) (ClubPlayerRow, error) {
	playerRow, err := newPlayerRow(myClubId, name)
	if err != nil {
		return ClubPlayerRow{}, fmt.Errorf("failed to create base player row: %w", err)
	}

	parsedRunPower, err := strconv.ParseFloat(runPower, 64)
	if err != nil {
		return ClubPlayerRow{}, fmt.Errorf("failed to parse run power: %w", err)
	}
	if parsedRunPower < 0.0 {
		return ClubPlayerRow{}, fmt.Errorf("run power too low: needs to be between 0-10")
	}
	if parsedRunPower > 10.0 {
		return ClubPlayerRow{}, fmt.Errorf("run power too high: needs to be between 0-10")
	}

	parsedBallHandling, err := strconv.ParseFloat(ballHandling, 64)
	if err != nil {
		return ClubPlayerRow{}, fmt.Errorf("failed to parse ball handling: %w", err)
	}

	if parsedBallHandling < 0.0 {
		return ClubPlayerRow{}, fmt.Errorf("ball handling too low: needs to be between 0-10")
	}
	if parsedBallHandling > 10.0 {
		return ClubPlayerRow{}, fmt.Errorf("ball handling too high: needs to be between 0-10")
	}

	return ClubPlayerRow{
		PlayerRow:    playerRow,
		RunPower:     parsedRunPower,
		BallHandling: parsedBallHandling,
	}, nil
}
