package player

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

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
