package main

import (
	"fmt"

	"github.com/Zigelzi/go-tiimit/db"
	"github.com/Zigelzi/go-tiimit/player"
	"github.com/Zigelzi/go-tiimit/practice"
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
	actions := []string{"Create teams for a practice", "Import players to club", "Exit"}
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
		practice.GetAttendees()
		err := practice.CreateTeams()
		if err != nil {
			fmt.Println(err)
			break
		}
		practice.PrintTeams()
	case actions[1]:
		err := player.ImportToClub("202412_Kuntofutis_Pelaajat.xlsx")
		if err != nil {
			fmt.Println(err)
		}
	case actions[len(actions)-1]:
		return false
	}
	return true
}
