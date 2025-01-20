package player_test

import (
	"testing"

	"github.com/Zigelzi/go-tiimit/player"
)

var tests = []struct {
	name         string
	actualPlayer player.Player
	expected     float64
}{
	{
		name:         "Player with high attributes",
		actualPlayer: player.New(123, "Matti Meikäläinen", 10, 10, false),
		expected:     22.0,
	},
	{
		name:         "Player with zero attributes",
		actualPlayer: player.New(123, "Matti Meikäläinen", 0, 0, false),
		expected:     0.0,
	},
	{
		name:         "Empty player struct",
		actualPlayer: player.Player{},
		expected:     0,
	},
	{
		name:         "Player with negative ball handling",
		actualPlayer: player.New(123, "Matti Meikäläinen", -10, 0, false),
		expected:     0,
	},
	{
		name:         "Player with negative run power",
		actualPlayer: player.New(123, "Matti Meikäläinen", 0, -10, false),
		expected:     0,
	},
	{
		name:         "Player with negative attributes",
		actualPlayer: player.New(123, "Matti Meikäläinen", -10, -10, false),
		expected:     0,
	},
}

func TestScore(t *testing.T) {
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			actual := testCase.actualPlayer.Score()
			if actual != testCase.expected {
				t.Errorf("Score() = %.2f: want %.2f", actual, testCase.expected)
			}
		})
	}
}
