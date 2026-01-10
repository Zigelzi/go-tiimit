package main

import "strings"

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
