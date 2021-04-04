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

			// Makes a record
			maxObtainedTkns := len(plr.TokensInRegionByStatus("NN", "NN", "A", 1))
			potentialCandidates := opponents[i].TokensInRegionByStatus("NN", "NN", "A", -1)
			opponents[i].RecordPotentialCandidates(maxObtainedTkns, potentialCandidates)
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

			// Makes a record
			maxObtainedTkns := len(plr.TokensInRegionByStatus("NN", "NN", "A", 1))
			potentialCandidates := opponents[i].TokensInRegionByStatus("NN", "NN", "A", -1)
			opponents[i].RecordPotentialCandidates(maxObtainedTkns, potentialCandidates)
		}
	}
}

// Number of tokens in a block
func NumTknsInRegion(start string, end string, terrain string) int {
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
	var xorSet []string

	if len(a) > len(b) {
		for _, va := range a {
			var found bool = false
			for _, vb := range b {
				if va == vb {
					found = true
				}
			}
			if !found {
				xorSet = append(xorSet, va)
			}
		}
		// If b is a's subset
		if len(a)-len(xorSet) == len(b) {
			return true, xorSet
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
				xorSet = append(xorSet, vb)
			}
		}
		// If a is b's subset
		if len(b)-len(xorSet) == len(a) {
			return true, xorSet
		}
	}

	return false, nil
}

func computeDeductedTokens(start string, end string, terrainType string, numReportedTkns int, tarNo string, opponents []player.Player) (int, []string) {
	var candidateTkns, obtainedTkns []string
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
		return -1, nil
	}

	for nTokens, potentialTknsList := range targetPOTList {
		if nTokens == numUnfoundObtainedTkns {
			for _, potentialTkns := range potentialTknsList {
				foundSubset, xorSet := isSubsetAndGetXOR(candidateTkns, potentialTkns)
				if foundSubset {
					return 0, xorSet
				}
			}
		} else if nTokens > numUnfoundObtainedTkns {
			for _, potentialTkns := range potentialTknsList {
				if len(potentialTkns)-len(candidateTkns) == nTokens-numUnfoundObtainedTkns {
					foundSubset, xorSet := isSubsetAndGetXOR(candidateTkns, potentialTkns)
					if foundSubset {
						return 1, xorSet
					}
				}
			}
		} else if nTokens < numUnfoundObtainedTkns {
			for _, potentialTkns := range potentialTknsList {
				if len(candidateTkns)-len(potentialTkns) == numUnfoundObtainedTkns-nTokens {
					foundSubset, xorSet := isSubsetAndGetXOR(candidateTkns, potentialTkns)
					if foundSubset {
						return 1, xorSet
					}
				}
			}
		}
	}

	return -1, nil
}

func computeTokensStatus(start string, end string, terrainType string, numReportedTkns int, tarNo string, opponents []player.Player) (int, []string) {
	var candidateTkns, obtainedTkns []string
	var numUnfoundObtainedTkns int

	// Finds the target player
	for _, opponent := range opponents {
		if tarNo == opponent.No {
			candidateTkns = opponent.TokensInRegionByStatus(start, end, terrainType, -1)
			obtainedTkns = opponent.TokensInRegionByStatus(start, end, terrainType, 1)
		}
	}
	numUnfoundObtainedTkns = numReportedTkns - len(obtainedTkns)
	// Case 1: Non of unfirmed tokens in the requested region or the checked token's status is firmed
	if len(candidateTkns) == 0 {
		return -1, nil
	} else if len(candidateTkns) >= 1 { // Case 2: 1 or more of unfirmed tokens in the requested region
		// Subcase 1
		if numUnfoundObtainedTkns == 0 {
			return 0, candidateTkns
		} else if numUnfoundObtainedTkns == len(candidateTkns) { // Subcase 2
			return 1, candidateTkns
		}
	}

	// Records potential obtained tokens
	for i, _ := range opponents {
		if tarNo == opponents[i].No {
			opponents[i].RecordPotentialCandidates(numUnfoundObtainedTkns, candidateTkns)
			// Printing Potential tokens report
			opponents[i].DisplayPotentialTokensReport()
		}
	}

	return -1, nil
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
	var numReportedTkns int
	var startPosition, endPosition, tarNo, terrain string

	startPosition = msg[0][:2]
	endPosition = msg[1][:2]
	tarNo = msg[4][1:]
	terrain = msg[2]
	numReportedTkns, _ = strconv.Atoi(msg[3])

	// Deducts potential tokens by comparing history report
	// Tokens in the set would be deducted from potential tokens for a firmed status
	dStatus, xorSet := computeDeductedTokens(startPosition, endPosition, terrain, numReportedTkns, tarNo, opponents)
	for _, dToken := range xorSet {
		updatePlayerTable(dStatus, dToken, tarNo, opponents, plr)
		// If the token is a treasure location
		if isTreasure(dToken, opponents, plr) {
			updatePlayerTable(2, dToken, tarNo, opponents, plr)
		}
	}

	// Computes current report of tokens' status
	// Tokens in the set would be setted for a firmed status
	status, tokenSet := computeTokensStatus(startPosition, endPosition, terrain, numReportedTkns, tarNo, opponents)
	for _, token := range tokenSet {
		updatePlayerTable(status, token, tarNo, opponents, plr)
		// If the token is a treasure location
		if isTreasure(token, opponents, plr) {
			updatePlayerTable(2, token, tarNo, opponents, plr)
		}
	}
}
