package player

func (player Player) Score() float64 {
	const runPowerWeight float64 = 1.2
	const ballHandlingWeight float64 = 1
	return player.runPower*runPowerWeight + player.ballHandling*ballHandlingWeight
}
