package player

func (player Player) Score() float64 {
	const runPowerWeight float64 = 1.2
	const ballHandlingWeight float64 = 1
	if player.ballHandling < 0 || player.runPower < 0 {
		return 0
	}
	return player.runPower*runPowerWeight + player.ballHandling*ballHandlingWeight
}
