package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Zigelzi/go-tiimit/internal/auth"
	"github.com/Zigelzi/go-tiimit/internal/db"
	"github.com/Zigelzi/go-tiimit/internal/file"
	"github.com/Zigelzi/go-tiimit/internal/player"
	"github.com/Zigelzi/go-tiimit/internal/practice"
	"github.com/manifoldco/promptui"
)

func main() {
	newDb, err := db.InitDB()
	if err != nil {
		log.Fatalf("initializing database failed: %v", err)
		return
	}
	defer newDb.Close()

	cfg := cliConfig{
		queries: db.New(newDb),
		db:      newDb,
	}
	for {
		if !selectAction(cfg) {
			break
		}
	}
}

func selectAction(cfg cliConfig) bool {
	// TODO: Move selecting create/import action to it's own function.
	actions := []string{
		"Create teams for a practice by importing MyClub attendees",
		"Manage players",
		"Create new user",
		"Exit",
	}
	prompt := promptui.Select{
		Label: "What do you want to do",
		Items: actions,
	}
	_, result, err := prompt.Run()
	if err != nil {
		fmt.Println("Unable to get input for selecting action")
		return false
	}

	switch result {
	case actions[0]:

		var attendanceDirectory = "attendance-files/"
		fileName, err := file.Select(attendanceDirectory)
		if err != nil {
			fmt.Println(err)
			break
		}

		attendancePlayerRows, _ := file.ImportAttendancePlayerRowsFromPath(attendanceDirectory + fileName)
		confirmedRows, err := file.GetAttendanceRowsByStatus(attendancePlayerRows, file.AttendanceIn)
		if err != nil {
			fmt.Println(err)
			break
		}

		dbConfirmedPlayers := []db.Player{}
		for _, row := range confirmedRows {
			confirmedDbPlayer, err := cfg.queries.GetPlayerByMyclubID(context.Background(), int64(row.PlayerRow.MyclubID))
			if err != nil {
				fmt.Println(err)
				continue
			}
			dbConfirmedPlayers = append(dbConfirmedPlayers, confirmedDbPlayer)
		}
		confirmedPlayers := []player.Player{}
		for _, dbConfirmedPlayer := range dbConfirmedPlayers {
			confirmedPlayers = append(confirmedPlayers, player.New(dbConfirmedPlayer.MyclubID, dbConfirmedPlayer.Name, dbConfirmedPlayer.RunPower, dbConfirmedPlayer.BallHandling, dbConfirmedPlayer.IsGoalie))
		}

		unknownRows, err := file.GetAttendanceRowsByStatus(attendancePlayerRows, file.AttendanceUnknown)
		if err != nil {
			fmt.Println(err)
			break
		}

		dbUnknownPlayers := []db.Player{}
		for _, row := range unknownRows {
			unknownDbPlayer, err := cfg.queries.GetPlayerByMyclubID(context.Background(), int64(row.PlayerRow.MyclubID))
			if err != nil {
				fmt.Println(err)
				continue
			}
			dbUnknownPlayers = append(dbUnknownPlayers, unknownDbPlayer)
		}

		unknownPlayers := []player.Player{}
		for _, dbUnknownPlayer := range dbUnknownPlayers {
			unknownPlayers = append(unknownPlayers, player.New(dbUnknownPlayer.MyclubID, dbUnknownPlayer.Name, dbUnknownPlayer.RunPower, dbUnknownPlayer.BallHandling, dbUnknownPlayer.IsGoalie))
		}

		player.SortByScore(confirmedPlayers)
		player.SortByScore(unknownPlayers)
		goalies, fieldPlayers := player.GetPreferences(confirmedPlayers)
		team1, team2, err := practice.Distribute(fieldPlayers, goalies)

		if err != nil {
			fmt.Println(err)
			break
		}

		newPractice := practice.Practice{
			TeamOnePlayers: team1,
			TeamTwoPlayers: team2,
			Date:           time.Now(),
		}

		newPractice.PrintTeams()
		for i, unknownPlayer := range unknownPlayers {
			fmt.Printf("%d. %s [%.1f]\n\n", i+1, unknownPlayer.Details(), unknownPlayer.Score())
		}

	case actions[1]:
		err := player.Manage(cfg.queries)
		if err != nil {
			fmt.Println(err)
		}
	case actions[2]:
		// Handler
		fmt.Println("Creating new user")
		fmt.Println("-----------------")

		fmt.Println("Write username of the new user:")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		username := cleanInput(scanner.Text())[0]

		fmt.Println("Write password of the new user:")
		scanner.Scan()
		password := cleanInput(scanner.Text())[0]

		trimmedPassword := strings.TrimSpace(password)
		if trimmedPassword == "" {
			fmt.Printf("password can't be empty\n")
			return false
		}

		if auth.IsWeakPassword(trimmedPassword) {
			fmt.Printf("password must be at least %d characters\n", auth.PasswordMinLength)
			return false
		}

		hashedPassword, err := auth.HashPassword(trimmedPassword)
		if err != nil {
			fmt.Printf("failed to hash password: %w\n", err)
			return false
		}

		newUser, err := cfg.queries.CreateUser(context.Background(), db.CreateUserParams{
			Username:       username,
			HashedPassword: hashedPassword,
		})
		if err != nil {
			fmt.Printf("failed to create new user: %v\n", err)
			return false
		}

		fmt.Printf("Created new user [%s] successfully!\n", newUser.Username)
	case actions[len(actions)-1]:
		// Exit should be always last action
		return false
	}
	return true
}
