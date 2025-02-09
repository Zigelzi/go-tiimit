package main

import (
	"fmt"

	"github.com/Zigelzi/go-tiimit/db"
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
		"Create teams for a practice manually",
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
		practice := practice.New()

		err := practice.MarkAttendees()
		if err != nil {
			fmt.Println(err)
			break
		}

		team1, team2, err := team.Distribute(practice.Players)
		if err != nil {
			fmt.Println(err)
			break
		}

		err = practice.AddTeams(team1, team2)
		if err != nil {
			fmt.Println(err)
			break
		}

		practice.PrintTeams()

	case actions[1]:
		practice := practice.New()

		players, err := player.ImportAttendees()
		if err != nil {
			fmt.Println(err)
			break
		}

		team1, team2, err := team.Distribute(players)
		if err != nil {
			fmt.Println(err)
			break
		}

		err = practice.AddTeams(team1, team2)
		if err != nil {
			fmt.Println(err)
			break
		}

		practice.PrintTeams()
	case actions[2]:
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
