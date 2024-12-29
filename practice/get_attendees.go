package practice

import (
	"fmt"

	"github.com/Zigelzi/go-tiimit/player"
	"github.com/fatih/color"
)

func (p *Practice) GetAttendees() error {
	players, err := player.GetAll()
	if err != nil {
		return fmt.Errorf("failed to load players %w", err)
	}

	fmt.Println("Mark which players are attending to create the teams.")
	fmt.Println("1 - Attends")
	fmt.Println("2 - Doesn't attend")
	fmt.Println("3 - Go to previous player")

	i := 0
AttendanceLoop:
	for {
		player := players[i]
		var selection string

		fmt.Printf("\n%s (%d/%d) \n", player.Name, i+1, len(players))
		fmt.Scanln(&selection)
		switch selection {
		case "1":
			err := p.Add(player)
			if err != nil {
				color.Red(err.Error())
				break
			}

			color.Green("Added player '%s' to attending players. Now %d players are attending", player.Name, len(p.Players))

			// Move to next unassigned player if it exists. If it doesn't it means that user has assigned all players.
			if i+1 < len(players) {
				i += 1
				continue
			}
			break AttendanceLoop

		case "2":
			isRemoved := p.Remove(player)

			if isRemoved {
				color.Yellow("Removed player '%s' from attending players. Now %d players are attending", player.Name, len(p.Players))
			} else {
				color.Yellow("Skipped player '%s'. Now %d players are attending", player.Name, len(p.Players))

			}

			// Move to next unassigned player if it exists. If it doesn't it means that user has assigned all players.
			if i+1 < len(players) {
				i += 1
				continue
			}
			break AttendanceLoop

		case "3":
			if i-1 >= 0 {
				i -= 1
			} else {
				color.Red("Can't go back. No previous player exists")
			}

		default:
			color.Red("No action for %s. Select action from the list.\n\n", selection)
		}

	}
	return nil
}
