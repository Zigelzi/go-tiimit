package player

import (
	"fmt"
	"strconv"

	"github.com/xuri/excelize/v2"
)

type Player struct {
	MyClubId     int64
	Name         string
	runPower     float64
	ballHandling float64
}

type ByScore []Player

func (p ByScore) Len() int           { return len(p) }
func (p ByScore) Less(i, j int) bool { return p[i].GetScore() > p[j].GetScore() }
func (p ByScore) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func New(id int64, name string, runPower float64, ballHandling float64) Player {
	return Player{
		MyClubId:     id,
		Name:         name,
		runPower:     runPower,
		ballHandling: ballHandling,
	}
}

func (player Player) GetScore() float64 {
	const runPowerWeight float64 = 1.2
	const ballHandlingWeight float64 = 1
	return player.runPower*runPowerWeight + player.ballHandling*ballHandlingWeight
}

func (player Player) PrintDetails() {
	fmt.Printf("[%d] %s score: %.1f\n", player.MyClubId, player.Name, player.GetScore())
}

func Load(fileName string) ([]Player, error) {
	file, err := excelize.OpenFile(fileName)
	if err != nil {
		return nil, err
	}
	defer closeFile(file)

	rows, err := file.GetRows("Tapahtuma")
	if err != nil {
		return nil, err
	}
	playerRows := rows[4:]
	var players []Player

	for _, playerRow := range playerRows {
		myClubId, err := strconv.Atoi(playerRow[0])
		if err != nil {
			fmt.Println("Unable to parse MyClub ID.")
			return nil, err
		}

		runPower, err := strconv.ParseFloat(playerRow[3], 64)
		if err != nil {
			fmt.Println("Unable to parse run power.")
			return nil, err
		}

		ballHandling, err := strconv.ParseFloat(playerRow[4], 64)
		if err != nil {
			fmt.Println("Unable to parse ball handling.")
			return nil, err
		}
		player := New(int64(myClubId), playerRow[1], runPower, ballHandling)
		players = append(players, player)
	}

	fmt.Printf("Loaded %d players from file %s\n", len(players), fileName)
	return players, nil
}

func closeFile(openFile *excelize.File) {
	if err := openFile.Close(); err != nil {
		fmt.Println(err)
		return
	}
}
