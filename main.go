package main

import (
	"errors"
	"fmt"

	"example.com/go-tiimit/team"
)

func main() {
	selectAction()
}

func selectAction() {
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

		team1.PrintPlayers()
		team2.PrintPlayers()
	case "9":
		return
	}
	selectAction()
}

func getAction() (action string, err error) {
	fmt.Println("What do you want to do?")
	fmt.Println("1 - Create teams for a practice")
	fmt.Println("9 - Exit")
	fmt.Scanln(&action)
	if action == "" {
		return "", errors.New("no input value provided")
	}

	return action, nil
}
