package tkncmp

import (
	"strconv"

	"../player"
)

var TokenMap = [24]string{"1F", "1B", "1M", "2F", "2B", "2M", "3F", "3B", "3M", "4F", "4B", "4M", "5F", "5B", "5M", "6F", "6B", "6M", "7F", "7B", "7M", "8F", "8B", "8M"}
var directionIndexMap = map[string]int{"NN": 0, "NE": 3, "EE": 6, "SE": 9, "SS": 12, "SW": 15, "WW": 18, "NW": 21}

func AllocatedTokensCompute(tokens []string, plr *player.Player, opponents []player.Player) {
	for _, token := range tokens {
		plr.MakeRecord(token, 1)
		for _, opponent := range opponents {
			opponent.MakeRecord(token, 0)
		}
	}
}

func LeftTokensCompute(tokens []string, plr *player.Player, opponents []player.Player) {
	for _, token := range tokens {
		plr.MakeRecord(token, 0)
		for _, opponent := range opponents {
			opponent.MakeRecord(token, 0)
		}
	}
}

func TokenInfoSwapCompute(token string, opponentNo string, plr *player.Player, opponents []player.Player) {
	plr.MakeRecord(token, 0)
	for _, opponent := range opponents {
		if opponent.No == opponentNo {
			opponent.MakeRecord(token, 1)
		} else {
			opponent.MakeRecord(token, 0)
		}
	}
}

func computeTokenStatus(start int, end int, terrainType string, reportedTrnNum int) int {
	// Certain-no
	if reportedTrnNum == 0 {
		return 0
	}
	// Certain-yes
	numTrnToLook := end - start
	if terrainType != "A" {
		numTrnToLook = numTrnToLook / 3
	}
	if reportedTrnNum == numTrnToLook {
		return 1
	}
	// Potential
	return 2
}

func PlayerReportCompute(msg []string, opponents []player.Player) {
	var idxFrom, idxEnd, countEnd, numTrn, status int
	var plrNo, terrain string

	idxFrom = directionIndexMap[msg[0][:2]]
	idxEnd = directionIndexMap[msg[1][:2]]
	countEnd = idxEnd
	// In the case ex: WW:18 to NE:3 or NN:0 to NN:0
	if idxEnd <= idxFrom {
		countEnd += 24
	}
	plrNo = msg[4][1:]
	terrain = msg[2]
	numTrn, _ = strconv.Atoi(msg[3])
	status = computeTokenStatus(idxFrom, countEnd, terrain, numTrn)

	for idxFrom != countEnd {
		token := TokenMap[idxFrom%24]
		if terrain == "A" || terrain == token[1:] {
			for _, opponent := range opponents {
				switch status {
				case 0:
					if plrNo == opponent.No {
						opponent.MakeRecord(token, 0)
					} else {
						opponent.MakeRecord(token, 2)
					}
				case 1:
					if plrNo == opponent.No {
						opponent.MakeRecord(token, 1)
					} else {
						opponent.MakeRecord(token, 0)
					}
				// Incomplete case, could find centern tokens
				case 2:
					if plrNo == opponent.No {
						opponent.MakeRecord(token, 2)
					}
				}
			}
		}
		idxFrom++
	}
}
