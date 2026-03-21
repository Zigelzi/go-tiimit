package player

import "math"

func (player Player) Score() float64 {
	const runPowerWeight float64 = 1.2
	const ballHandlingWeight float64 = 1
	if player.ballHandling < 0 || player.runPower < 0 {
		return 0
	}
	weightedScore := player.runPower*runPowerWeight + player.ballHandling*ballHandlingWeight
	return roundToTwoDecimals(weightedScore)
}

func roundToTwoDecimals(value float64) float64 {
	return math.Round(value*100) / 100
}
