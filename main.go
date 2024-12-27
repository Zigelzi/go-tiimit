package main

import (
	"errors"
	"fmt"

	"github.com/Zigelzi/go-tiimit/db"
	"github.com/Zigelzi/go-tiimit/practice"
)

func main() {
	db.Init()
	db.CreateTables()
	defer db.DB.Close()
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
		practice := practice.New()
		practice.GetAttendees()
		err := practice.CreateTeams()
		if err != nil {
			fmt.Println(err)
			break
		}
		practice.PrintTeams()

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
