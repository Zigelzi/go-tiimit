package player

import (
	"reflect"
	"testing"
)

var goalies = []Player{
	{Name: "Miika", IsGoalie: true},
	{Name: "Pekka", IsGoalie: true},
	{Name: "Jouni", IsGoalie: true},
	{Name: "Pia", IsGoalie: true},
	{Name: "Seppo", IsGoalie: true},
}

var fieldPlayers = []Player{
	{Name: "Janne", IsGoalie: false},
	{Name: "Antti", IsGoalie: false},
	{Name: "Anne", IsGoalie: false},
	{Name: "Mauno", IsGoalie: false},
	{Name: "Kirsi", IsGoalie: false},
}

var tests = []struct {
	name                 string
	initialPlayers       []Player
	expectedGoalies      []Player
	expectedFieldPlayers []Player
}{
	{
		name:                 "No players",
		initialPlayers:       []Player{},
		expectedGoalies:      []Player{},
		expectedFieldPlayers: []Player{},
	},
	{
		name:                 "Only goalies",
		initialPlayers:       goalies,
		expectedGoalies:      goalies,
		expectedFieldPlayers: []Player{},
	},
	{
		name:                 "Only field players",
		initialPlayers:       fieldPlayers,
		expectedGoalies:      []Player{},
		expectedFieldPlayers: fieldPlayers,
	},
	{
		name:                 "5 goalies and 5 field players",
		initialPlayers:       append(append([]Player{}, goalies...), fieldPlayers...),
		expectedGoalies:      goalies,
		expectedFieldPlayers: fieldPlayers,
	},
	{
		name:                 "1 goalie and 5 field players",
		initialPlayers:       append(append([]Player{}, goalies[:1]...), fieldPlayers...),
		expectedGoalies:      goalies[:1],
		expectedFieldPlayers: fieldPlayers,
	},
	{
		name:                 "5 goalies and 1 field players",
		initialPlayers:       append(append([]Player{}, goalies...), fieldPlayers[:1]...),
		expectedGoalies:      goalies,
		expectedFieldPlayers: fieldPlayers[:1],
	},
	{
		name:                 "nil players",
		initialPlayers:       nil,
		expectedGoalies:      []Player{},
		expectedFieldPlayers: []Player{},
	},
}

func TestGetPreferences(t *testing.T) {
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			actualGoalies, actualFieldPlayers := GetPreferences(testCase.initialPlayers)

			if actualGoalies == nil {
				t.Fatalf("expected goalies to be non-nil, got nil")
			}
			if actualFieldPlayers == nil {
				t.Fatalf("expected field players to be non-nil, got nil")
			}

			if len(actualGoalies) != len(testCase.expectedGoalies) {
				t.Errorf("number of goalies doesn't match: got %d want %d", len(actualGoalies), len(testCase.expectedGoalies))
			}
			if len(actualFieldPlayers) != len(testCase.expectedFieldPlayers) {
				t.Errorf("number of goalies doesn't match: got %d want %d", len(actualFieldPlayers), len(testCase.expectedFieldPlayers))
			}

			if !reflect.DeepEqual(actualGoalies, testCase.expectedGoalies) {
				t.Errorf("goalies don't match: got %v want %v", actualGoalies, testCase.expectedGoalies)
			}
		})
	}
}
