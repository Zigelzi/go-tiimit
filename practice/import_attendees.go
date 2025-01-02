package practice

import (
	"fmt"
	"strconv"

	"github.com/Zigelzi/go-tiimit/file"
	"github.com/Zigelzi/go-tiimit/player"
	"github.com/xuri/excelize/v2"
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
	file, err := excelize.OpenFile(attendanceDirectory + fileName)
	if err != nil {
		return fmt.Errorf("unable to open file to import attendees from a file %s: %w", fileName, err)
	}
	defer closeFile(file)

	rows, err := file.GetRows("Tapahtuma")
	if err != nil {
		return fmt.Errorf("unable to read rows to import attendees from a file: %w", err)
	}

	// List of players in MyClub start on row 5. Rows before that are other details or empty.
	playerRows := rows[4:]
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

	fmt.Printf("Imported %d players from file %s\n", len(addedPlayers), fileName)
	return nil
}

func closeFile(openFile *excelize.File) {
	if err := openFile.Close(); err != nil {
		fmt.Println(err)
		return
	}
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
