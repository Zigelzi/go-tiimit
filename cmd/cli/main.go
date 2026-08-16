package main

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Zigelzi/go-tiimit/internal/auth"
	"github.com/Zigelzi/go-tiimit/internal/db"
	"github.com/joho/godotenv"
	"github.com/manifoldco/promptui"
)

type cliConfig struct {
	queries *db.Queries
	db      *sql.DB
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("failed to load .env file: %v", err)
	}
	newDb, err := db.InitDB()
	if err != nil {
		log.Fatalf("initializing database failed: %v", err)
		return
	}
	defer newDb.Close()

	err = db.RunMigrations(newDb)
	if err != nil {
		log.Fatalf("failed to run migrations: %v", err)
		return
	}

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
	actions := []string{
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
			fmt.Printf("failed to hash password: %v\n", err)
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

func cleanInput(text string) []string {
	text = strings.Trim(text, " ")
	words := strings.Split(text, " ")
	cleanedWords := make([]string, len(words))

	for i, word := range words {
		lowerCaseWord := strings.ToLower(word)
		cleanedWords[i] = lowerCaseWord
	}
	return cleanedWords
}
