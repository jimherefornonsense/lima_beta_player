package player

import (
	"fmt"
	"strconv"
)

// Status flag -1: no mark, 0: unobtained, 1: obtained, 2: potential obtained 3: treasure
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

var TableIndexMap = map[string]int{"NN": 0, "NE": 1, "EE": 2, "SE": 3, "SS": 4, "SW": 5, "WW": 6, "NW": 7}

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
// Status flag -1: no mark, 0: unobtained, 1: obtained, 2: potential obtained 3: treasure
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

// Checks token status
func (plr *Player) StatusByToken(token string) int {
	var status int = -1

	region, _ := strconv.Atoi(token[0:1])
	terrain := token[1:]
	if terrain == "B" {
		status = plr.Table[region-1].Beach
	} else if terrain == "F" {
		status = plr.Table[region-1].Forest
	} else if terrain == "M" {
		status = plr.Table[region-1].Mountain
	}

	return status
}

//Reports tokens in a block according to its status
func (plr *Player) TokensInRegionByStatus(start string, end string, terrain string, checkedStatus int) []string {
	var itStart, itEnd int
	var tokens []string

	itStart = TableIndexMap[start]
	itEnd = TableIndexMap[end]
	if itEnd <= itStart {
		itEnd += 8
	}
	for itStart != itEnd {
		it := itStart % 8
		if terrain == "B" {
			if plr.Table[it].Beach == checkedStatus {
				tokens = append(tokens, strconv.Itoa(it+1)+"B")
			}
		} else if terrain == "F" {
			if plr.Table[it].Forest == checkedStatus {
				tokens = append(tokens, strconv.Itoa(it+1)+"F")
			}
		} else if terrain == "M" {
			if plr.Table[it].Mountain == checkedStatus {
				tokens = append(tokens, strconv.Itoa(it+1)+"M")
			}
		} else if terrain == "A" {
			if plr.Table[it].Beach == checkedStatus {
				tokens = append(tokens, strconv.Itoa(it+1)+"B")
			}
			if plr.Table[it].Forest == checkedStatus {
				tokens = append(tokens, strconv.Itoa(it+1)+"F")
			}
			if plr.Table[it].Mountain == checkedStatus {
				tokens = append(tokens, strconv.Itoa(it+1)+"M")
			}
		}
		itStart++
	}

	return tokens
}

// Checks tokens in a block of status -1 and 2
func (plr *Player) UnfirmedTokensInRegion(start string, end string, terrain string) []string {
	return append(plr.TokensInRegionByStatus(start, end, terrain, -1), plr.TokensInRegionByStatus(start, end, terrain, 2)...)
}

// Checks tokens in a block of status -1
func (plr *Player) UnfirmedOneTokensInRegion(start string, end string, terrain string) []string {
	return append(plr.TokensInRegionByStatus(start, end, terrain, -1))
}

// Checks tokens in a block of status 2
func (plr *Player) UnfirmedTwoTokensInRegion(start string, end string, terrain string) []string {
	return append(plr.TokensInRegionByStatus(start, end, terrain, 2))
}



// Prints the matrix of current table
func (plr *Player) DisplayTable() {
	fmt.Printf("--------------Player%s--------------\n", plr.No)
	fmt.Printf("%10sNN NE EE SE SS SW WW NW NN\n", "Direction:")
	fmt.Printf("%11s", "Region: ")
	for i := 1; i < 9; i++ {
		fmt.Printf("R%d ", i)
	}
	fmt.Print("\n")
	fmt.Printf("%11s", "Beach: ")
	for _, trs := range plr.Table {
		fmt.Printf("%2d ", trs.Beach)
	}
	fmt.Print("\n")
	fmt.Printf("%11s", "Forest: ")
	for _, trs := range plr.Table {
		fmt.Printf("%2d ", trs.Forest)
	}
	fmt.Print("\n")
	fmt.Printf("%11s", "Mountain: ")
	for _, trs := range plr.Table {
		fmt.Printf("%2d ", trs.Mountain)
	}
	fmt.Print("\n")
}
