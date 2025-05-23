package main

import (
	"fmt"

	"github.com/Zigelzi/go-tiimit/db"
	"github.com/Zigelzi/go-tiimit/file"
	"github.com/Zigelzi/go-tiimit/player"
	"github.com/Zigelzi/go-tiimit/practice"
	"github.com/Zigelzi/go-tiimit/team"
	"github.com/manifoldco/promptui"
)

func main() {
	db.Init()
	db.CreateTables()
	defer db.DB.Close()
	for {
		if !selectAction() {
			break
		}
	}
}

func selectAction() bool {
	// TODO: Move selecting create/import action to it's own function.
	actions := []string{
		"Create teams for a practice by importing MyClub attendees",
		"Manage players",
		"Exit",
	}
	prompt := promptui.Select{
		Label: "What do you want to do",
		Items: actions,
	}
	_, result, err := prompt.Run()
	if err != nil {
		fmt.Println("Unable to get input for selecting action")
		return false
	}

	switch result {
	case actions[0]:
		newPractice := practice.New()

		var attendanceDirectory = "attendance-files/"
		fileName, err := file.Select(attendanceDirectory)
		if err != nil {
			fmt.Println(err)
			break
		}

		attendancePlayerRows, _ := file.ImportAttendancePlayerRows(attendanceDirectory + fileName)
		for _, row := range attendancePlayerRows {
			err := newPractice.AddPlayer(row.PlayerRow.MyClubId, row.Attendance)
			if err != nil {
				fmt.Println(err)
				continue
			}
		}
		confirmedPlayers, err := newPractice.GetPlayersByStatus(practice.AttendanceIn, player.Get)
		if err != nil {
			fmt.Println(err)
		}
		unknownPlayers, err := newPractice.GetPlayersByStatus(practice.AttendanceUnknown, player.Get)
		if err != nil {
			fmt.Println(err)
		}

		player.SortByScore(confirmedPlayers)
		player.SortByScore(unknownPlayers)
		goalies, fieldPlayers := player.GetPreferences(confirmedPlayers)
		team1, team2, err := team.Distribute(goalies, fieldPlayers)

		if err != nil {
			fmt.Println(err)
			break
		}

		err = newPractice.AddTeams(team1, team2)
		if err != nil {
			fmt.Println(err)
			break
		}

		newPractice.PrintTeams()
		for i, unknownPlayer := range unknownPlayers {
			fmt.Printf("%d. %s [%.1f]\n\n", i+1, unknownPlayer.Details(), unknownPlayer.Score())
		}

	case actions[1]:
		err := player.Manage()
		if err != nil {
			fmt.Println(err)
		}
	case actions[len(actions)-1]:
		// Exit should be always last action
		return false
	}
	return true
}
