package player

import (
	"fmt"
	"strconv"
)

// Status flag -1: no mark, 0: certain-no, 1: certain-yes, 2: potential
type terrains struct {
	Beach    int
	Forest   int
	Mountain int
}

// Player struct
// Table's indices corresponse to island map's regions, ex: index 0 = region 1, index 1 = region 2
type Player struct {
	No             string
	PlayerTerrains []string
	Table          [8]terrains
}

func NewPlayer(playerNo string) Player {
	plr := Player{No: playerNo}
	for i := range plr.Table {
		plr.Table[i].Beach = -1
		plr.Table[i].Forest = -1
		plr.Table[i].Mountain = -1
	}
	return plr
}

// Parses each token and makes a record
// Status flag -1: no mark, 0: certain-no, 1: certain-yes, 2: potential
// token format ex: 1B, 2F, 3M
func (plr *Player) MakeRecord(token string, tokenStatus int) {
	region, _ := strconv.Atoi(token[0:1])
	terrain := token[1:]
	if terrain == "B" {
		if plr.Table[region-1].Beach == -1 || plr.Table[region-1].Beach == 2 {
			plr.Table[region-1].Beach = tokenStatus
		}
	} else if terrain == "F" {
		if plr.Table[region-1].Forest == -1 || plr.Table[region-1].Forest == 2 {
			plr.Table[region-1].Forest = tokenStatus
		}
	} else if terrain == "M" {
		if plr.Table[region-1].Mountain == -1 || plr.Table[region-1].Mountain == 2 {
			plr.Table[region-1].Mountain = tokenStatus
		}
	}
}

// Prints the matrix of current table
func (plr *Player) DisplayTable() {
	fmt.Printf("-------------Player%s-------------\n", plr.No)
	fmt.Printf("%10s", "Region: ")
	for i := 1; i < 9; i++ {
		fmt.Printf("R%d ", i)
	}
	fmt.Print("\n")
	fmt.Printf("%10s", "Beach: ")
	for _, trs := range plr.Table {
		fmt.Printf("%2d ", trs.Beach)
	}
	fmt.Print("\n")
	fmt.Printf("%10s", "Forest: ")
	for _, trs := range plr.Table {
		fmt.Printf("%2d ", trs.Forest)
	}
	fmt.Print("\n")
	fmt.Printf("%10s", "Mountain: ")
	for _, trs := range plr.Table {
		fmt.Printf("%2d ", trs.Mountain)
	}
	fmt.Print("\n")
}
