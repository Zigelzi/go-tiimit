package main

import (
	"errors"
	"fmt"

	"example.com/go-tiimit/db"
	"example.com/go-tiimit/team"
)

func main() {
	db.Init()
	selectAction()
}

func selectAction() {
	fmt.Println("What do you want to do?")
	fmt.Println("1 - Create teams for a practice")
	fmt.Println("9 - Exit")
	action, err := getAction()
	if err != nil {
		fmt.Println(err)
	}
	switch action {
	case "1":

		team1, team2, err := team.CreatePracticeTeams()

		if err != nil {
			fmt.Println(err)
			break
		}

		team1.Details()
		team2.Details()
	case "9":
		return
	default:
		fmt.Printf("No action for %s. Select action from the list.\n\n", action)
	}
	selectAction()
}

func getAction() (action string, err error) {

	fmt.Scanln(&action)
	if action == "" {
		return "", errors.New("no input value provided")
	}

	return action, nil
}
