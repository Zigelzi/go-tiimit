package player

import (
	"context"
	"fmt"

	"github.com/Zigelzi/go-tiimit/internal/db"
	"github.com/manifoldco/promptui"
)

// Manage lets user select and execute player management actions.
func Manage(dbQuery *db.Queries) error {
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
		err := ImportToClub(dbQuery)
		if err != nil {
			fmt.Println(err)
		}
	case actions[1]:
		dbPlayers, err := dbQuery.GetAllPlayers(context.Background())
		if err != nil {
			return err
		}
		players := []Player{}
		for _, dbPlayer := range dbPlayers {
			players = append(players, FromDB(dbPlayer))
		}
		chosenPlayer, err := choose("Select player to edit goalie status of", players)
		if err != nil {
			return err
		}
		err = dbQuery.ToggleGoalieStatus(context.Background(), db.ToggleGoalieStatusParams{
			ID:       chosenPlayer.ID,
			IsGoalie: !chosenPlayer.IsGoalie,
		})
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
func choose(label string, players []Player) (Player, error) {
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
