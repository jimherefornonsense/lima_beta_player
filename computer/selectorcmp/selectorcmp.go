package selectorcmp

import (
	"fmt"
	"lima_beta_player/player"
	"sort"
	"strings"
)

type Combination struct {
	No           string
	Region       [2]int
	Terrain      string
	Potentials   []string
	Unchosen_die int
}

// Checks there's no the same set in history sets
func isNewPotentialSet(candidates []string, opponent player.Player) bool {
	if len(candidates) == 0 {
		return false
	}

	for _, sets := range opponent.PotentialObtainedTknsList {
		for _, set := range sets {
			if len(set) == len(candidates) {
				var equalityCounter int = 0
				setString := fmt.Sprintf("%s", set)

				for _, token := range candidates {
					if strings.Contains(setString, token) {
						equalityCounter++
					}
				}

				// Evaluating
				if equalityCounter == len(candidates) {
					// A same set in record
					return false
				}
			}
		}
	}

	return true
}

func terrainParser(t1 string, t2 string) []string {
	var terrains = []string{}

	if t1 == "W" && t2 == "W" {
		terrains = append(terrains, "A")
		terrains = append(terrains, "B")
		terrains = append(terrains, "F")
		terrains = append(terrains, "M")
	} else if t1 == "W" && t2 != "W" {
		terrains = append(terrains, "A")
		terrains = append(terrains, t2)
	} else if t1 != "W" && t2 == "W" {
		terrains = append(terrains, "A")
		terrains = append(terrains, t1)
	} else if t1 == t2 {
		terrains = append(terrains, t1)
	} else {
		terrains = append(terrains, "A")
	}

	return terrains
}

// Returns true if found a Combination intelligently, false if not
func Selection(rolledDice []string, opponents []player.Player) (bool, Combination) {
	var combinations_group []Combination
	var temp Combination

	for j := range opponents {
		for i := 1; i <= 3; i++ {
			for k := i + 1; k <= 3; k++ {
				terrains := terrainParser(rolledDice[i][2:], rolledDice[k][2:])
				for _, t := range terrains {
					var c1 Combination
					c1.No = opponents[j].No
					c1.Region = [2]int{i, k}
					c1.Terrain = t
					c1.Potentials = opponents[j].UnfirmedOneTokensInRegion(rolledDice[i][:2], rolledDice[k][:2], t)
					for l := 1; l <= 3; l++ {
						if l != i && l != k {
							c1.Unchosen_die = l
						}
					}

					var c2 Combination
					c2.No = opponents[j].No
					c2.Region = [2]int{k, i}
					c2.Terrain = t
					c2.Potentials = opponents[j].UnfirmedOneTokensInRegion(rolledDice[k][:2], rolledDice[i][:2], t)
					for l := 1; l <= 3; l++ {
						if l != i && l != k {
							c2.Unchosen_die = l
						}
					}

					if isNewPotentialSet(c1.Potentials, opponents[j]) {
						combinations_group = append(combinations_group, c1)
					}
					if isNewPotentialSet(c2.Potentials, opponents[j]) {
						combinations_group = append(combinations_group, c2)
					} else { // Just in case there is no new potential set or is an empty set
						temp.No = c2.No
						temp.Region[0] = c2.Region[0]
						temp.Region[1] = c2.Region[1]
						temp.Terrain = c2.Terrain
						temp.Unchosen_die = c2.Unchosen_die
					}
				}
			}
		}
	}

	// Sorting in descending order
	sort.SliceStable(combinations_group, func(i, j int) bool {
		return len(combinations_group[i].Potentials) > len(combinations_group[j].Potentials)
	})

	if len(combinations_group) != 0 {
		return true, combinations_group[0]
	}

	return false, temp
}

// If true, using pistol, and returns a target number
func IsUsePistol(plr player.Player, opponents []player.Player) (bool, string) {
	var notATarget int = -1

	for i := range opponents {
		if opponents[i].Pistol < 1 {
			notATarget = i
		}
	}

	if notATarget != -1 {
		var targetI int
		if notATarget+1 == len(opponents) {
			notATarget = 0
		} else {
			targetI = notATarget + 1
		}
		return true, opponents[targetI].No
	} else if len(plr.UnfirmedOneTokensInRegion("NN", "NN", "A")) < 5 {
		return true, opponents[len(opponents)-1].No
	}

	return false, ""
}

func Pistoling(rolledDice []string, targetNo string, opponents []player.Player) (string, string, string, string) {
	var makeUpMsg = []string{"P0", "NNW", "EEW", "SWW"}
	_, aCombination := Selection(makeUpMsg, opponents)
	d1, d2 := rolledDice[aCombination.Region[0]], rolledDice[aCombination.Region[1]]

	return d1, d2, aCombination.Terrain, targetNo
}

func Shoveling(die string) string {
	var terrainSet = [3]string{"B", "F", "M"}
	var terrain string

	for _, t := range terrainSet {
		if die[2:] != t {
			terrain = t
			break
		}
	}

	return terrain
}

func Barreling(rolledDice []string, aCombination Combination) []int {
	var rerolled_dice []int

	rerolled_dice = append(rerolled_dice, aCombination.Region[0]-1)
	rerolled_dice = append(rerolled_dice, aCombination.Region[1]-1)
	if rolledDice[aCombination.Unchosen_die][2:] != "W" {
		rerolled_dice = append(rerolled_dice, aCombination.Unchosen_die-1)
	}

	return rerolled_dice
}
