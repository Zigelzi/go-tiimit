package player

import (
	"fmt"
	"strconv"

	"github.com/Zigelzi/go-tiimit/file"
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

func ImportAttendees() (players []Player, err error) {
	var importedPlayers []Player
	fileName, err := file.Select(attendanceDirectory)
	if err != nil {
		fmt.Println(err)
		return importedPlayers, err
	}

	playerRows, err := file.ImportRows(attendanceDirectory + fileName)
	if err != nil {
		return importedPlayers, err
	}

	for _, playerRow := range playerRows {
		status := attendanceName[playerRow[3]]
		if status == AttendanceIn {
			player, err := getPlayer(playerRow[0])
			if err != nil {
				return importedPlayers, fmt.Errorf("unable to get player: %w", err)
			}
			importedPlayers = append(importedPlayers, player)
		}
	}

	fmt.Printf("Imported %d players to practice from a file %s\n\n", len(importedPlayers), fileName)
	return importedPlayers, nil
}

func getPlayer(rowContent string) (Player, error) {
	myClubId, err := strconv.Atoi(rowContent)
	if err != nil {
		return Player{}, fmt.Errorf("unable to parse MyClub ID %s: %w", rowContent, err)
	}

	newPlayer, err := Get(int64(myClubId))
	if err != nil {
		return Player{}, fmt.Errorf("unable to add attending player with MyClubID %s: %w", rowContent, err)
	}
	return newPlayer, nil
}
