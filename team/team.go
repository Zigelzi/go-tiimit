package team

import (
	"errors"
	"fmt"
	"sort"
	"strconv"

	"example.com/go-tiimit/player"
	"github.com/xuri/excelize/v2"
)

type Team struct {
	name    string
	players []player.Player
}

func CreatePracticeTeams() (team1 Team, team2 Team, err error) {
	attendingPlayers, err := getAttendees()
	if err != nil {
		return Team{}, Team{}, err
	}

	team1.name = "Team 1"
	team2.name = "Team 2"
	if len(attendingPlayers) == 0 {
		return Team{}, Team{}, errors.New("no attending players to distribute")
	}

	sort.Sort(player.ByScore(attendingPlayers))

	for i, player := range attendingPlayers {
		if (i+1)&2 == 0 {
			team1.players = append(team1.players, player)
		} else {
			team2.players = append(team2.players, player)
		}
	}

	return team1, team2, nil
}

func (team *Team) Details() {
	fmt.Printf("%s players\n", team.name)
	for i, player := range team.players {
		fmt.Printf("%d. %s\n", i+1, player.Name)
	}
	fmt.Printf("\n%s has %d players with total score of %.1f\n\n", team.name, len(team.players), team.score())
}

func (team *Team) score() float64 {
	totalScore := 0.0
	for _, player := range team.players {
		totalScore += player.GetScore()
	}
	return totalScore
}

func getAttendees() ([]player.Player, error) {
	var attendingPlayers []player.Player

	fmt.Println("How do you want to mark the attending players?")
	fmt.Println("1 - Mark attendance manually")
	// fmt.Println("2 - Load attendees from a file")

	var attendanceInput string
	fmt.Scanln(&attendanceInput)
	switch attendanceInput {
	case "1":
		attendingPlayers, err := markAttendeesManually()
		if err != nil {
			fmt.Println("Failed to mark attendees manually")
			return nil, err
		}
		return attendingPlayers, err
	default:
		fmt.Printf("No action for %s. Select action from the list\n\n", attendanceInput)
		getAttendees()
	}

	return attendingPlayers, nil
}

func markAttendeesManually() ([]player.Player, error) {
	team, err := loadTeamFromFile("202412_Kuntofutis_Pelaajat.xlsx")
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("failed to load players from a file")
	}

	fmt.Println("Mark which players are attending to create the teams.")
	fmt.Println("1 - Attends")
	fmt.Println("2 - Doesn't attend")
	fmt.Println("3 - Go to previous player")

	var attendingPlayers []player.Player

	i := 0
AttendanceLoop:
	for {
		player := team.players[i]
		var selection string

		fmt.Printf("%s (%d/%d) \n", player.Name, i+1, len(team.players))
		fmt.Scanln(&selection)
		switch selection {
		case "1":
			attendingPlayers = append(attendingPlayers, player)
			if i+1 < len(team.players) {
				i += 1
			} else {
				break AttendanceLoop
			}
		case "2":
			if i+1 < len(team.players) {
				i += 1
			} else {
				break AttendanceLoop
			}
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
	return attendingPlayers, nil
}

func loadTeamFromFile(fileName string) (*Team, error) {
	file, err := excelize.OpenFile(fileName)
	if err != nil {
		return nil, err
	}
	defer closeFile(file)

	rows, err := file.GetRows("Tapahtuma")
	if err != nil {
		return nil, err
	}
	playerRows := rows[4:]
	var newTeam Team

	for _, playerRow := range playerRows {
		myClubId, err := strconv.Atoi(playerRow[0])
		if err != nil {
			fmt.Println("Unable to parse MyClub ID.")
			return nil, err
		}

		runPower, err := strconv.ParseFloat(playerRow[3], 64)
		if err != nil {
			fmt.Println("Unable to parse run power.")
			return nil, err
		}

		ballHandling, err := strconv.ParseFloat(playerRow[4], 64)
		if err != nil {
			fmt.Println("Unable to parse ball handling.")
			return nil, err
		}
		player := player.New(int64(myClubId), playerRow[1], runPower, ballHandling)
		newTeam.players = append(newTeam.players, player)
	}

	fmt.Printf("Loaded %d players from file %s\n", len(newTeam.players), fileName)
	return &newTeam, nil
}

func closeFile(openFile *excelize.File) {
	if err := openFile.Close(); err != nil {
		fmt.Println(err)
		return
	}
}
