package practice

import (
	"fmt"
	"strings"

	"github.com/Zigelzi/go-tiimit/player"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
)

func (p *Practice) MarkAttendees() error {
	players, err := player.GetAll()
	if err != nil {
		return fmt.Errorf("failed to load players %w", err)
	}

	playerIndex := 0
	actions := []string{"Attends", "Doesn't attend", "Go to previous player"}
	prompt := promptui.Select{
		Items:        actions,
		HideSelected: true,
	}

	printInstructions()
AttendanceLoop:
	for {
		player := players[playerIndex]
		prompt.Label = fmt.Sprintf("%s (%d/%d)", player.Name, playerIndex+1, len(players))

		_, result, err := prompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			break
		}

		switch result {
		case actions[0]:

			if p.Add(player) {
				color.Green("Added player '%s' to attending players. Now %d players are attending", player.Name, len(p.Players))
			} else {
				color.Yellow("Skipped '%s' as they're already attending. Now %d players are attending", player.Name, len(p.Players))
			}

			// Move to next unassigned player if it exists. If it doesn't it means that user has assigned all players.
			if playerIndex+1 < len(players) {
				playerIndex += 1
				continue
			}
			break AttendanceLoop

		case actions[1]:

			if p.Remove(player) {
				color.Yellow("Removed player '%s' from attending players. Now %d players are attending", player.Name, len(p.Players))
			} else {
				color.Yellow("Skipped player '%s'. Now %d players are attending", player.Name, len(p.Players))

			}

			// Move to next unassigned player if it exists. If it doesn't it means that user has assigned all players.
			if playerIndex+1 < len(players) {
				playerIndex += 1
				continue
			}
			break AttendanceLoop

		case actions[2]:
			if playerIndex-1 >= 0 {
				playerIndex -= 1
			} else {
				color.Red("Can't go back. No previous player exists")
			}
		}

	}
	return nil
}

func printInstructions() {
	instruction := "Mark players attendance to the practice"

	fmt.Println(strings.Repeat("=", len(instruction)))
	fmt.Println(instruction)
	fmt.Println(strings.Repeat("=", len(instruction)))
}
