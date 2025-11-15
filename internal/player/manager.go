package player

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

// Manage lets user select and execute player management actions.
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
		chosenPlayer, err := choose("Select player to edit goalie status of")
		if err != nil {
			return err
		}
		err = ToggleGoalieStatus(chosenPlayer)
		if err != nil {
			return err
		}
	case actions[len(actions)-1]:
		return nil
	}
	return nil
}

// Choose a player from all existing players.
// Takes [label] as input to display the desired text for customizing the action.
func choose(label string) (Player, error) {
	players, err := GetAll()
	if err != nil {
		return Player{}, err
	}
	templates := &promptui.SelectTemplates{
		Inactive: "  {{ .Details }}",
		Active:   fmt.Sprintf("%s {{ .Details | underline }}", promptui.IconSelect),
		Selected: fmt.Sprintf(`{{ "%s" | green }} {{ .Details | faint }}`, promptui.IconGood),
	}

	prompt := promptui.Select{
		Label:     label,
		Items:     players,
		Templates: templates,
	}
	i, _, err := prompt.Run()
	if err != nil {
		return Player{}, err
	}

	chosenPlayer := players[i]
	return chosenPlayer, nil
}
