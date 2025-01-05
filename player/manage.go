package player

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

func Manage() error {
	actions := []string{
		"Import players to club",
		"Update goalie status",
		"Go back",
	}
	prompt := promptui.Select{
		Label: "What do you want to do",
		Items: actions,
	}
	_, result, err := prompt.Run()
	if err != nil {
		fmt.Println("Unable to get input for selecting action when managing players")
		return err
	}

	switch result {
	case actions[0]:
		err := ImportToClub()
		if err != nil {
			fmt.Println(err)
		}
	case actions[1]:
		fmt.Println("Player goalie status updated!")
	case actions[len(actions)-1]:
		return nil
	}
	return nil
}
