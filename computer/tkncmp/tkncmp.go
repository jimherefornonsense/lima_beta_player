package tkncmp

import (
	"strconv"

	"lima_beta_player/player"
)

var TokenMap = [24]string{"1B", "1F", "1M", "2B", "2F", "2M", "3B", "3F", "3M", "4B", "4F", "4M", "5B", "5F", "5M", "6B", "6F", "6M", "7B", "7F", "7M", "8B", "8F", "8M"}
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

func isSubsetAndGetXOR(a []string, b []string) (bool, []string) {
	var xORset []string

	if len(a) > len(b) {
		for _, va := range a {
			var found bool = false
			for _, vb := range b {
				if va == vb {
					found = true
				}
			}
			if !found {
				xORset = append(xORset, va)
			}
		}
		// If b is a's subset
		if len(a)-len(xORset) == len(b) {
			return true, xORset
		}
	} else {
		for _, vb := range b {
			var found bool = false
			for _, va := range a {
				if vb == va {
					found = true
				}
			}
			if !found {
				xORset = append(xORset, vb)
			}
		}
		// If a is b's subset
		if len(b)-len(xORset) == len(a) {
			return true, xORset
		}
	}

	return false, xORset
}

func computeDeductedTokens(start string, end string, terrainType string, numReportedTkns int, tarNo string, opponents []player.Player) (int, []string) {
	var candidateTkns, obtainedTkns, xORset []string
	var targetPOTList map[int][][]string
	var numUnfoundObtainedTkns int

	// Finds the target player
	for _, opponent := range opponents {
		if tarNo == opponent.No {
			candidateTkns = opponent.TokensInRegionByStatus(start, end, terrainType, -1)
			obtainedTkns = opponent.TokensInRegionByStatus(start, end, terrainType, 1)
			targetPOTList = opponent.PotentialObtainedTknsList
		}
	}
	numUnfoundObtainedTkns = numReportedTkns - len(obtainedTkns)
	if len(candidateTkns) < 1 {
		return -1, xORset
	}

	for nTokens, potentialTknsList := range targetPOTList {
		if nTokens == numUnfoundObtainedTkns {
			for _, potentialTkns := range potentialTknsList {
				foundSubset, xORset := isSubsetAndGetXOR(candidateTkns, potentialTkns)
				if foundSubset {
					return 0, xORset
				}
			}
		} else if nTokens > numUnfoundObtainedTkns {
			for _, potentialTkns := range potentialTknsList {
				if len(potentialTkns)-len(candidateTkns) == nTokens-numUnfoundObtainedTkns {
					foundSubset, xORset := isSubsetAndGetXOR(candidateTkns, potentialTkns)
					if foundSubset {
						return 1, xORset
					}
				}
			}
		} else if nTokens < numUnfoundObtainedTkns {
			for _, potentialTkns := range potentialTknsList {
				if len(candidateTkns)-len(potentialTkns) == numUnfoundObtainedTkns-nTokens {
					foundSubset, xORset := isSubsetAndGetXOR(candidateTkns, potentialTkns)
					if foundSubset {
						return 1, xORset
					}
				}
			}
		}
	}

	return -1, xORset
}

func computeTokenStatus(start string, end string, terrainType string, numReportedTkns int, token string, tarNo string, opponents []player.Player) int {
	var candidateTkns, obtainedTkns []string
	var numUnfoundObtainedTkns, curTknStatus int

	// Finds the target player
	for _, opponent := range opponents {
		if tarNo == opponent.No {
			candidateTkns = opponent.TokensInRegionByStatus(start, end, terrainType, -1)
			obtainedTkns = opponent.TokensInRegionByStatus(start, end, terrainType, 1)
			curTknStatus = opponent.StatusByToken(token)
		}
	}
	numUnfoundObtainedTkns = numReportedTkns - len(obtainedTkns)
	// Case 1: Non of unfirmed tokens in the requested region or the checked token's status is firmed
	if len(candidateTkns) == 0 || curTknStatus != -1 {
		return -1
	} else if len(candidateTkns) == 1 { // Case 2: 1 of unfirmed tokens in the requested region
		// Subcase 1
		if numUnfoundObtainedTkns == 0 {
			return 0
		} else if numUnfoundObtainedTkns == 1 { // Subcase 2
			return 1
		}
	} else if len(candidateTkns) > 1 { // Case 3: More than 1 of unfirmed tokens in the requested region
		// Subcase 1
		if numUnfoundObtainedTkns == 0 {
			return 0
		} else if numUnfoundObtainedTkns == len(candidateTkns) { // Subcase 2
			return 1
		}
	}
	// Records potential obtained tokens
	if token == candidateTkns[len(candidateTkns)-1] {
		for i, _ := range opponents {
			if tarNo == opponents[i].No {
				opponents[i].RecordPotentialCandidates(numUnfoundObtainedTkns, candidateTkns)
			}
		}
	}

	return -1
}

func updatePlayerTable(status int, token string, tarNo string, opponents []player.Player, plr *player.Player) {
	for i, opponent := range opponents {
		switch status {
		case 0: // The target player doesn't have the token
			if tarNo == opponent.No {
				opponents[i].MakeRecord(token, 0)
			}
		case 1: // The target player has the token
			if tarNo == opponent.No {
				opponents[i].MakeRecord(token, 1)
				plr.MakeRecord(token, 0)
			} else {
				opponents[i].MakeRecord(token, 0)
			}
		case 2: // Found the treasure
			if tarNo == opponent.No {
				opponents[i].MakeRecord(token, 0)
				plr.MakeRecord(token, 2)
			} else {
				opponents[i].MakeRecord(token, 0)
			}
		default: // Does nothing
		}
	}
}

func isTreasure(token string, opponents []player.Player, plr *player.Player) bool {
	var found bool = true

	if plr.StatusByToken(token) != -1 {
		found = false
	}
	for _, opponent := range opponents {
		if opponent.StatusByToken(token) != 0 {
			found = false
		}
	}

	return found
}

func PlayerReportCompute(msg []string, plr *player.Player, opponents []player.Player) {
	var idxFrom, idxEnd, itEnd, numReportedTkns int
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
	numReportedTkns, _ = strconv.Atoi(msg[3])

	// Deducts potential tokens by comparing history report
	// Tokens in the set would be deducted from potential tokens for a firmed status
	dStatus, xORset := computeDeductedTokens(msg[0][:2], msg[1][:2], terrain, numReportedTkns, tarNo, opponents)
	for _, dToken := range xORset {
		updatePlayerTable(dStatus, dToken, tarNo, opponents, plr)
		// If the token is a treasure location
		if isTreasure(dToken, opponents, plr) {
			updatePlayerTable(2, dToken, tarNo, opponents, plr)
		}
	}

	for idxFrom != itEnd {
		// A current token in the searched region
		token := TokenMap[idxFrom%24]
		if terrain == "A" || terrain == token[1:] {
			status := computeTokenStatus(msg[0][:2], msg[1][:2], terrain, numReportedTkns, token, tarNo, opponents)
			updatePlayerTable(status, token, tarNo, opponents, plr)
			// If the token is a treasure location
			if isTreasure(token, opponents, plr) {
				updatePlayerTable(2, token, tarNo, opponents, plr)
			}
		}
		idxFrom++
	}
}
