package practice

import (
	"fmt"
	"strconv"

	"github.com/Zigelzi/go-tiimit/file"
	"github.com/Zigelzi/go-tiimit/player"
)

type AttendanceStatus int

const (
	AttendanceIn AttendanceStatus = iota
	AttendanceOut
	AttendanceUnknown
)

const attendanceDirectory = "attendance-files/"

var attendanceName = map[string]AttendanceStatus{
	"Osallistuu":   AttendanceIn,
	"Ei osallistu": AttendanceOut,
	"Ei vastausta": AttendanceUnknown,
}

func (p *Practice) ImportAttendees() error {
	fileName, err := file.Select(attendanceDirectory)
	if err != nil {
		fmt.Println(err)
		return err
	}

	playerRows, err := file.ImportRows(attendanceDirectory + fileName)
	if err != nil {
		return err
	}

	var addedPlayers []player.Player

	for _, playerRow := range playerRows {
		status := attendanceName[playerRow[3]]
		if status == AttendanceIn {
			player, err := getPlayer(playerRow[0])
			if err != nil {
				return fmt.Errorf("unable to get player: %w", err)
			}
			p.Add(player)
			addedPlayers = append(addedPlayers, player)
		}
	}

	fmt.Printf("Imported %d players to practice from a file %s\n\n", len(addedPlayers), fileName)
	return nil
}

func getPlayer(rowContent string) (player.Player, error) {
	myClubId, err := strconv.Atoi(rowContent)
	if err != nil {
		return player.Player{}, fmt.Errorf("unable to parse MyClub ID %s: %w", rowContent, err)
	}

	newPlayer, err := player.Get(int64(myClubId))
	if err != nil {
		return player.Player{}, fmt.Errorf("unable to add attending player with MyClubID %s: %w", rowContent, err)
	}
	return newPlayer, nil
}
