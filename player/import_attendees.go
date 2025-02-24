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

func ImportAttendees() (confirmedPlayers, unknownPlayers []Player, err error) {
	confirmedPlayers = []Player{}
	unknownPlayers = []Player{}

	// Import player rows
	fileName, err := file.Select(attendanceDirectory)
	if err != nil {
		fmt.Println(err)
		return confirmedPlayers, unknownPlayers, err
	}

	playerRows, err := file.ImportPlayerRows(attendanceDirectory + fileName)
	if err != nil {
		return confirmedPlayers, unknownPlayers, err
	}

	// Get players who could attend
	for _, playerRow := range playerRows {
		status := attendanceName[playerRow[3]]
		player, err := getPlayer(playerRow[0])
		if err != nil {
			return confirmedPlayers, unknownPlayers, fmt.Errorf("unable to get player: %w", err)
		}
		if status == AttendanceIn {
			confirmedPlayers = append(confirmedPlayers, player)
		}
		if status == AttendanceUnknown {
			unknownPlayers = append(unknownPlayers, player)
		}
	}

	fmt.Printf("Imported %d players to practice from a file %s\n\n", len(confirmedPlayers)+len(unknownPlayers), fileName)
	return confirmedPlayers, unknownPlayers, nil
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
