package player_test

import (
	"testing"

	"github.com/Zigelzi/go-tiimit/player"
)

var detailsTests = []struct {
	name     string
	player   player.Player
	expected string
}{
	{
		name:     "Player is goalie",
		player:   player.Player{Name: "Matti Meikäläinen", IsGoalie: true},
		expected: "Matti Meikäläinen [G]",
	},
	{
		name:     "Player is not a goalie",
		player:   player.Player{Name: "Matti Meikäläinen", IsGoalie: false},
		expected: "Matti Meikäläinen",
	},
	{
		name:     "Player goalie status is not known",
		player:   player.Player{Name: "Matti Meikäläinen"},
		expected: "Matti Meikäläinen",
	},
	{
		name:     "Player name is empty and is not goalie",
		player:   player.Player{Name: "", IsGoalie: false},
		expected: "",
	},
	{
		name:     "Player name is empty and is goalie",
		player:   player.Player{Name: "", IsGoalie: true},
		expected: " [G]",
	},
}

func TestDetails(t *testing.T) {
	for _, testCase := range detailsTests {
		t.Run(testCase.name, func(t *testing.T) {
			actual := testCase.player.Details()
			if actual != testCase.expected {
				t.Errorf("Details() = %q: want %q", actual, testCase.expected)
			}
		})
	}
}
