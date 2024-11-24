package team

import (
	"errors"
	"fmt"
	"strconv"

	"example.com/go-tiimit/player"
	"github.com/xuri/excelize/v2"
)

type Team struct {
	name    string
	players []player.Player
}

func CreatePracticeTeams() (team1 Team, team2 Team, err error) {
	attendingPlayers, err := markAttendance()
	team1.name = "Team 1"
	team2.name = "Team 2"
	if err != nil {
		return Team{}, Team{}, err
	}
	if len(attendingPlayers) == 0 {
		return Team{}, Team{}, errors.New("no attending players to distribute")
	}

	for i, player := range attendingPlayers {
		if (i+1)&2 == 0 {
			team1.players = append(team1.players, player)
		} else {
			team2.players = append(team2.players, player)
		}
	}

	return team1, team2, nil
}

func (team *Team) PrintPlayers() {
	fmt.Printf("%s players\n", team.name)
	for _, player := range team.players {
		player.PrintDetails()
	}
	fmt.Printf("%s has %d players\n\n", team.name, len(team.players))
}

func markAttendance() ([]player.Player, error) {
	team, err := loadTeamFromFile()

	fmt.Println("Mark which players are attending to create the teams.")
	fmt.Println("1 - Attends")
	fmt.Println("2 - Doesn't attend")

	var attendingPlayers []player.Player
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("failed to load players from a file")
	}

	for i, player := range team.players {
		var selection string

		fmt.Printf("%s (%d/%d) \n", player.Name, i+1, len(team.players))
		fmt.Scanln(&selection)
		if selection == "1" {
			attendingPlayers = append(attendingPlayers, player)
		}
	}

	return attendingPlayers, nil
}

func loadTeamFromFile() (*Team, error) {
	fileName := "Kuntofutis_Pelaajat.xlsx"
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
		playerId, err := strconv.Atoi(playerRow[3])
		if err != nil {
			return nil, err
		}
		player := player.New(int64(playerId), playerRow[1])
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
