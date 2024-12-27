package player

import (
	"fmt"
	"strconv"

	"github.com/xuri/excelize/v2"
)

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

		err = Insert(player)
		if err != nil {
			return nil, err
		}
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
