package file

import (
	"fmt"
	"strconv"
)

type PlayerRow struct {
	MyClubId int
	Name     string
}

func newPlayerRow(newMyClubId, name string) (PlayerRow, error) {
	myClubId, err := strconv.Atoi(newMyClubId)
	if err != nil {
		return PlayerRow{}, fmt.Errorf("unable to convert MyClubId to integer: %w", err)
	}

	if name == "" {
		return PlayerRow{}, fmt.Errorf("player name can't be empty")
	}

	return PlayerRow{
			MyClubId: myClubId,
			Name:     name,
		},
		nil
}
