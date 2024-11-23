package main

import (
	"errors"
	"fmt"
	"strconv"

	"example.com/go-tiimit/player"
	"github.com/xuri/excelize/v2"
)

func main() {

	players, err := loadPlayersFromFile()

	if err != nil {
		fmt.Println(err)
		return
	}
	selectAction(players)
}

func loadPlayersFromFile() ([]player.Player, error) {
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
	var players []player.Player

	for _, playerRow := range playerRows {
		playerId, err := strconv.Atoi(playerRow[3])
		if err != nil {
			return nil, err
		}
		player := player.New(int64(playerId), playerRow[1])
		players = append(players, player)
	}

	fmt.Printf("Loaded %d players from file %s\n", len(players), fileName)
	return players, nil
}

func closeFile(openFile *excelize.File) {
	if err := openFile.Close(); err != nil {
		fmt.Println(err)
		return
	}
}

func selectAction(players []player.Player) {
	action, err := getAction()
	if err != nil {
		fmt.Println(err)
	}
	switch action {
	case "1":
		attendingPlayers := markAttendance(players)
		team1, team2, err := createPracticeTeams(attendingPlayers)

		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("Attending today: %d\n\n", len(attendingPlayers))

		fmt.Println("Team 1 players")
		for _, player := range team1 {
			player.PrintDetails()
		}
		fmt.Printf("%d players in team 1\n\n", len(team1))

		fmt.Println("Team 2 players")

		for _, player := range team2 {
			player.PrintDetails()
		}
		fmt.Printf("Total players in team 2: %d\n", len(team2))
	case "9":
		return
	}
	selectAction(players)
}

func getAction() (action string, err error) {
	fmt.Println("What do you want to do?")
	fmt.Println("1 - Mark attending players")
	fmt.Println("9 - Exit")
	fmt.Scanln(&action)
	if action == "" {
		return "", errors.New("no input value provided")
	}

	return action, nil
}

func markAttendance(players []player.Player) []player.Player {
	var attendingPlayers []player.Player
	fmt.Println("Mark which players are attending to create the teams.")
	fmt.Println("1 - Attends")
	fmt.Println("2 - Doesn't attend")
	for i, player := range players {
		fmt.Printf("%s (%d/%d) \n", player.Name, i+1, len(players))
		var selection string
		fmt.Scanln(&selection)
		if selection == "1" {
			attendingPlayers = append(attendingPlayers, player)
		}
	}

	return attendingPlayers
}

func createPracticeTeams(attendingPlayers []player.Player) (team1 []player.Player, team2 []player.Player, err error) {
	if len(attendingPlayers) == 0 {
		return nil, nil, errors.New("no attending players to distribute")
	}

	for i, player := range attendingPlayers {
		if (i+1)&2 == 0 {
			team1 = append(team1, player)
		} else {
			team2 = append(team2, player)
		}
	}

	return team1, team2, nil
}
