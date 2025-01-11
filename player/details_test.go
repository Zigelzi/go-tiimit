package player_test

import (
	"testing"

	"github.com/Zigelzi/go-tiimit/player"
)

func TestDetails(t *testing.T) {
	tests := []struct {
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
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			actual := testCase.player.Details()
			if actual != testCase.expected {
				t.Errorf("Details() = %q: want %q", actual, testCase.expected)
			}
		})
	}
}
