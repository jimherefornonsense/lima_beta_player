package tkncmp

import (
	"strconv"

	"lima_beta_player/player"
)

var TokenMap = [24]string{"1F", "1B", "1M", "2F", "2B", "2M", "3F", "3B", "3M", "4F", "4B", "4M", "5F", "5B", "5M", "6F", "6B", "6M", "7F", "7B", "7M", "8F", "8B", "8M"}
var directionIndexMap = map[string]int{"NN": 0, "NE": 3, "EE": 6, "SE": 9, "SS": 12, "SW": 15, "WW": 18, "NW": 21}

func AllocatedTokensCompute(tokens []string, plr *player.Player, opponents []player.Player) {
	for _, token := range tokens {
		plr.MakeRecord(token, 1)
		for i := range opponents {
			opponents[i].MakeRecord(token, 0)
		}
	}
}

func LeftTokensCompute(tokens []string, plr *player.Player, opponents []player.Player) {
	for _, token := range tokens {
		plr.MakeRecord(token, 0)
		for i := range opponents {
			opponents[i].MakeRecord(token, 0)
		}
	}
}

func TokenInfoSwapCompute(token string, opponentNo string, plr *player.Player, opponents []player.Player) {
	plr.MakeRecord(token, 0)
	for i, opponent := range opponents {
		if opponent.No == opponentNo {
			opponents[i].MakeRecord(token, 1)
		} else {
			opponents[i].MakeRecord(token, 0)
		}
	}
}

// Number of tokens in a block
func numTknsInRegion(start string, end string, terrain string) int {
	var idxFrom, idxEnd, nToken int

	idxFrom = directionIndexMap[start]
	idxEnd = directionIndexMap[end]
	if idxEnd <= idxFrom {
		idxEnd += 24
	}
	nToken = idxEnd - idxFrom

	if terrain != "A" {
		return nToken / 3
	}
	return nToken
}

func computeTokenStatus(start string, end string, terrainType string, reportedTknNum int, tarNo string, opponents []player.Player) int {
	var tknsToLook, optainedTkns []string

	for _, opponent := range opponents {
		if tarNo == opponent.No {
			tknsToLook = opponent.UnfirmedTokensInRegion(start, end, terrainType)
			optainedTkns = opponent.TokensInRegionByStatus(start, end, terrainType, 1)
		}
	}
	// Unobtained / no unfirmed tokens
	if reportedTknNum == 0 || len(tknsToLook) == 0 {
		return 0
	}
	// Obtained
	if reportedTknNum-len(optainedTkns) == len(tknsToLook) {
		return 1
	} else if reportedTknNum-len(optainedTkns) < len(tknsToLook) {
		// Treasure check when only one token for checking
		if len(tknsToLook) == 1 {
			var isTreasure bool = true
			for _, opponent := range opponents {
				if tarNo != opponent.No {
					if opponent.StatusByToken(tknsToLook[0]) != 0 {
						isTreasure = false
					}
				}
			}
			if isTreasure {
				return 3
			} else { // Unobtained
				return 0
			}
		}
	}
	// Potential obtained
	return 2
}

func PlayerReportCompute(msg []string, plr *player.Player, opponents []player.Player) {
	var idxFrom, idxEnd, itEnd, reportedTknNum, status int
	var tarNo, terrain string

	idxFrom = directionIndexMap[msg[0][:2]]
	idxEnd = directionIndexMap[msg[1][:2]]
	itEnd = idxEnd
	// In the case ex: WW:18 to NE:3 or NN:0 to NN:0
	if idxEnd <= idxFrom {
		itEnd += 24
	}
	tarNo = msg[4][1:]
	terrain = msg[2]
	reportedTknNum, _ = strconv.Atoi(msg[3])
	status = computeTokenStatus(msg[0][:2], msg[1][:2], terrain, reportedTknNum, tarNo, opponents)

	for idxFrom != itEnd {
		token := TokenMap[idxFrom%24]
		if terrain == "A" || terrain == token[1:] {
			for i, opponent := range opponents {
				switch status {
				case 0:
					if tarNo == opponent.No {
						opponents[i].MakeRecord(token, 0)
					} else {
						opponents[i].MakeRecord(token, 2)
					}
				case 1:
					if tarNo == opponent.No {
						opponents[i].MakeRecord(token, 1)
						plr.MakeRecord(token, 0)
					} else {
						opponents[i].MakeRecord(token, 0)
					}
				case 2:
					if tarNo == opponent.No {
						opponents[i].MakeRecord(token, 2)
					}
				case 3:
					if tarNo == opponent.No {
						opponents[i].MakeRecord(token, 0)
						plr.MakeRecord(token, 3)
					} else {
						opponents[i].MakeRecord(token, 0)
					}
				}
			}
		}
		idxFrom++
	}
}
