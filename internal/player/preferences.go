package player

func GetPreferences(players []Player) (goalies, fieldPlayers []Player) {
	goalies = []Player{}
	fieldPlayers = []Player{}
	for _, player := range players {
		if player.IsGoalie {
			goalies = append(goalies, player)
		} else {
			fieldPlayers = append(fieldPlayers, player)
		}
	}

	return goalies, fieldPlayers
}
