package practice

import (
	"errors"
	"fmt"

	"github.com/Zigelzi/go-tiimit/player"
)

func (p *Practice) GetAttendees() error {
	players, err := player.Load("202412_Kuntofutis_Pelaajat.xlsx")
	if err != nil {
		fmt.Println(err)
		return errors.New("failed to load players from a file")
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

		fmt.Printf("%s (%d/%d) \n", player.Name, i+1, len(players))
		fmt.Scanln(&selection)
		switch selection {
		case "1":
			err := p.Add(player)
			if err != nil {
				fmt.Println(err)
				break
			}

			if i+1 < len(players) {
				i += 1
				continue
			}
			break AttendanceLoop

		case "2":
			err := p.Remove(player)
			if err != nil {
				fmt.Println(err)
				break
			}
			if i+1 < len(players) {
				i += 1
				continue
			}
			break AttendanceLoop

		case "3":
			if i-1 >= 0 {
				i -= 1
			} else {
				fmt.Println("Can't go back. No previous player exists")
			}

		default:
			fmt.Printf("No action for %s. Select action from the list.\n\n", selection)
		}

	}
	return nil
}
