package main

import (
	"fmt"

	"github.com/Zigelzi/go-tiimit/internal/db"
	"github.com/Zigelzi/go-tiimit/internal/file"
	"github.com/Zigelzi/go-tiimit/internal/player"
	"github.com/Zigelzi/go-tiimit/internal/practice"
	"github.com/Zigelzi/go-tiimit/internal/team"
	"github.com/manifoldco/promptui"
)

func main() {
	db.Init()
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

		attendancePlayerRows, _ := file.ImportAttendancePlayerRowsFromPath(attendanceDirectory + fileName)
		for _, row := range attendancePlayerRows {
			err := newPractice.AddPlayer(row.PlayerRow.MyClubId, row.Attendance)
			if err != nil {
				fmt.Println(err)
				continue
			}
		}

		dbConfirmedPlayers, _, err := newPractice.GetPlayersByStatus(practice.AttendanceIn, player.Get)
		if err != nil {
			fmt.Println(err)
		}
		confirmedPlayers := []player.Player{}
		for _, dbConfirmedPlayer := range dbConfirmedPlayers {
			confirmedPlayers = append(confirmedPlayers, player.New(dbConfirmedPlayer.MyClubId, dbConfirmedPlayer.Name, dbConfirmedPlayer.RunPower, dbConfirmedPlayer.BallHandling, dbConfirmedPlayer.IsGoalie))
		}

		dbUnknownPlayers, _, err := newPractice.GetPlayersByStatus(practice.AttendanceUnknown, player.Get)
		if err != nil {
			fmt.Println(err)
		}

		unknownPlayers := []player.Player{}
		for _, dbUnknownPlayer := range dbUnknownPlayers {
			unknownPlayers = append(unknownPlayers, player.New(dbUnknownPlayer.MyClubId, dbUnknownPlayer.Name, dbUnknownPlayer.RunPower, dbUnknownPlayer.BallHandling, dbUnknownPlayer.IsGoalie))
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
