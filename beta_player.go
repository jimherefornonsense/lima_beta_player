package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"lima_beta_player/computer/tkncmp"
	"lima_beta_player/human"
	"lima_beta_player/player"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
)

// Game struct
type Game struct {
	totalPlayers int
	autopilot    bool
	leftTokens   []string
}

type Combination struct {
	No         string
	regions    [2]int
	terrain    string
	potentials []string
}

type map_entry struct {
	val int
	key string
}

type map_entries []map_entry

func (s map_entries) Len() int           { return len(s) }
func (s map_entries) Less(i, j int) bool { return s[i].val < s[j].val }
func (s map_entries) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

var g Game
var p player.Player
var opponents []player.Player

type fn func(string) string

func selectedFunction(f fn, val string) string { // selectedFunction provides functionality to call specific function by its id [:2] of args string
	return f(val)
}

var functions = map[string]fn{
	"01": g.playerNO,
	"02": readMyTerrain,
	"03": g.leftoverTokens,
	"04": playerTurn,
	//"05": chooseDice,
	"06": interrogationReport,
	//"07": guessTokens,
	"08": guessCorrect,
	"09": guessIncorrect,
	"10": tokenInfoSwap,
	"11": remainingWinner,
	"14": pistolReport,
	"15": interrogationReport,
	"99": errorMsg,
}

func errorMsg(arg string) string {
	fmt.Println(arg)

	return ""
}

func (gm *Game) playerNO(args string) string {
	gm.totalPlayers, _ = strconv.Atoi(args[len(args)-1:])
	fmt.Printf("There are %d players, and you are player %s.\n", gm.totalPlayers, p.No)
	initOpponents(gm.totalPlayers)

	return ""
}

func initOpponents(totalPlayers int) {
	for i := 1; i <= totalPlayers; i++ {
		if strconv.Itoa(i) != p.No {
			plr := player.NewPlayer(strconv.Itoa(i))
			opponents = append(opponents, plr)
		}
	}
}

func readMyTerrain(args string) string {
	playerTokens := strings.Split(args[6:], ",")

	tkncmp.AllocatedTokensCompute(playerTokens, &p, opponents)

	// Init lists of potential obtained tokens for opponets
	for i, _ := range opponents {
		opponents[i].InitPotentialObtainedTknsList(len(playerTokens))
	}

	fmt.Println("Terrains are: " + strings.Join(playerTokens, ", "))

	return ""
}

func (gm *Game) leftoverTokens(args string) string {
	gm.leftTokens = strings.Split(args[3:], ",")

	tkncmp.LeftTokensCompute(gm.leftTokens, &p, opponents)

	fmt.Printf("Leftover tokens: ")
	for _, x := range gm.leftTokens {
		fmt.Printf(x + " ")
	}
	fmt.Printf("\n")

	return ""
}

func tokenInfoSwap(args string) string {
	message := strings.Split(args[3:], ",")

	if string(message[0][1:]) == p.No {
		fmt.Printf("You let %s know you got a token %s\n", message[1], message[2])
	} else {
		tkncmp.TokenInfoSwapCompute(message[2], message[0][1:], &p, opponents)

		fmt.Printf("You acknowledge %s got a token %s\n", message[0], message[2])
	}

	return ""
}

func remainingWinner(args string) string {
	message := strings.Split(args[3:], ",")

	fmt.Printf("%s wins as the only remaining player. All others have guessed incorrectly and been disqualified. The treasures are located at %s and %s\n",
		message[0], message[1], message[2])
	os.Exit(0)

	return ""
}

func playerTurn(args string) string {
	message := strings.Split(args[3:], ",")

	fmt.Println("Player " + message[0][1:] + " has rolled " + message[1] + "," + message[2] + "," + message[3])
	if "P"+p.No != message[0] {
		return ""
	}

	if g.autopilot {
		p.DisplayTable()
		for _, opponent := range opponents {
			opponent.DisplayTable()
		}

		isGuessing, tokens := p.IsGuessingAndGetAnswer()

		if isGuessing {
			return guessTokens(tokens[0], tokens[1])
		}

	} else {
		p.DisplayTable()
		for _, opponent := range opponents {
			opponent.DisplayTable()
		}
		if human.IsGuessing() {
			return guessTokens(human.GuessTreasures())
		}
	}

	return chooseDice(args)
}

// Special abilities
func pistol(d1 string, d2 string, t string, pNo string) string {
	var msg string

	msg = "05:" + d1 + "W," + d2 + "W," + t + "," + "P" + pNo + "," + "P"

	return msg
}

func shovel(selectedDice []string, index int, swappedTerrain string) {
	selectedDice[index] = strings.Replace(selectedDice[index], selectedDice[index][2:], swappedTerrain, 1)
	fmt.Println(selectedDice[index])
}

func terrainParser_comp(t1 string, t2 string, plr player.Player) []string {
	// var t string = "A"

	// var m = make(map[string]int)
	// m["A"] = len(plr.UnfirmedOneTokensInRegion(t1[:2], t2[:2], "A"))
	// m["B"] = len(plr.UnfirmedOneTokensInRegion(t1[:2], t2[:2], "B"))
	// m["F"] = len(plr.UnfirmedOneTokensInRegion(t1[:2], t2[:2], "F"))
	// m["M"] = len(plr.UnfirmedOneTokensInRegion(t1[:2], t2[:2], "M"))

	// var es map_entries
	// for k, v := range m {
	// 	es = append(es, map_entry{val: v, key: k})
	// }

	// sort.Sort(sort.Reverse(es))

	// if t1[2:] == "W" && t2[2:] == "W" {
	// 	t = es[0].key
	// 	//t = "A"
	// 	return t
	// } else if t1[2:] == "W" && t2[2:] != "W" {
	// 	t = "A"
	// } else if t1[2:] != "W" && t2[2:] == "W" {
	// 	fmt.Println("third:")
	// 	t = "A"
	// } else if t1[2:] == t2[2:] {
	// 	fmt.Println("fourth:")
	// 	return t1[2:]
	// }
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
	}

	return terrains
}

// Checks there's no the same set in history sets
func isNewPotentialSet(candidates []string, plr player.Player) bool {
	for _, sets := range plr.PotentialObtainedTknsList {
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

func chooseDice(args string) string {
	rolledDice := strings.Split(args[3:], ",")
	var n int
	var die1 int = 0
	var die2 int = 0
	var terrain, plr string

	// Temporary using random choosing, still implementing intelligent one
	if g.autopilot {
		var combinations_group []Combination
		// opponents_sets := make(map[string][]string)

		for j := range opponents {
			for i := 1; i <= 3; i++ {
				for k := i + 1; k <= 3; k++ {
					terrains := terrainParser_comp(rolledDice[i][2:], rolledDice[k][2:], opponents[j])
					for _, t := range terrains {
						var c1 Combination
						c1.No = opponents[j].No
						c1.regions = [2]int{i, k}
						c1.terrain = t
						c1.potentials = opponents[j].UnfirmedOneTokensInRegion(rolledDice[i][:2], rolledDice[k][:2], t)

						var c2 Combination
						c2.No = opponents[j].No
						c2.regions = [2]int{k, i}
						c2.terrain = t
						c2.potentials = opponents[j].UnfirmedOneTokensInRegion(rolledDice[k][:2], rolledDice[i][:2], t)

						if isNewPotentialSet(c1.potentials, opponents[j]) {
							combinations_group = append(combinations_group, c1)
						}
						if isNewPotentialSet(c2.potentials, opponents[j]) {
							combinations_group = append(combinations_group, c2)
						}

						// Just in case there is no new potential set
						plr = c1.No
						die1 = c1.regions[0]
						die2 = c1.regions[1]
						terrain = c1.terrain
					}
					// var c1 Combinations

					// c1.terrain = terrainParser_comp(rolledDice[i], rolledDice[k], opponents[j])
					// c1.potentials = opponents[j].UnfirmedOneTokensInRegion(rolledDice[i][:2], rolledDice[k][:2], c1.terrain)

					// c1.num_tokens = tkncmp.NumTknsInRegion(rolledDice[i][:2], rolledDice[k][:2], c1.terrain)
					// c1.regions = [2]int{i, k}
					// c1.No = opponents[j].No

					// var c2 Combinations
					// c2.terrain = terrainParser_comp(rolledDice[k], rolledDice[i], opponents[j])
					// c2.potentials = opponents[j].UnfirmedOneTokensInRegion(rolledDice[k][:2], rolledDice[i][:2], c2.terrain)

					// c2.num_tokens = tkncmp.NumTknsInRegion(rolledDice[k][:2], rolledDice[i][:2], c2.terrain)
					// c2.regions = [2]int{k, i}
					// c2.No = opponents[j].No

					// combinations_group = append(combinations_group, c1)
					// combinations_group = append(combinations_group, c2)
				}
			}

			// var temp []string
			// for l := 1; l <= len(opponents[j].PotentialObtainedTknsList); l++ {
			// 	for _, set := range opponents[j].PotentialObtainedTknsList[l] {
			// 		setString := fmt.Sprintf("%s", set)
			// 		temp = append(temp, setString)
			// 	}
			// }
			// opponents_sets[opponents[j].No] = temp
		}

		sort.SliceStable(combinations_group, func(i, j int) bool {
			return len(combinations_group[i].potentials) > len(combinations_group[j].potentials)
		})

		// for i := 0; i < len(combinations_group); i++ {
		// 	_, found := Find(opponents_sets[combinations_group[i].No], fmt.Sprintf("%s", combinations_group[i].potentials))
		// 	if !found {
		// 		fmt.Println("check")
		// 		plr = combinations_group[i].No
		// 		die1 = combinations_group[i].regions[0]
		// 		die2 = combinations_group[i].regions[1]
		// 		terrain = combinations_group[i].terrain
		// 		break
		// 	}

		// }
		if len(combinations_group) != 0 {
			plr = combinations_group[0].No
			die1 = combinations_group[0].regions[0]
			die2 = combinations_group[0].regions[1]
			terrain = combinations_group[0].terrain
		}

		// if die1 == 0 || die2 == 0 {
		// 	sort.SliceStable(combinations_group, func(i, j int) bool {
		// 		return combinations_group[i].num_tokens < combinations_group[i].num_tokens
		// 	})

		// 	plr = combinations_group[0].No
		// 	die1 = combinations_group[0].regions[0]
		// 	die2 = combinations_group[0].regions[1]
		// 	terrain = combinations_group[0].terrain
		// 	fmt.Println("heyyy")
		// }
		// fmt.Println(die1)
		// fmt.Println(die2)
		// fmt.Println(plr)

		var temp string = "05:" + rolledDice[die1] + "," + rolledDice[die2] + "," + terrain + ",P" + plr

		// fmt.Println(len(opponents))
		fmt.Println("temp-" + temp)

		return temp
	}

	// Manual
	spA := human.IsUsingSpA(&p)
	if spA == "P" {
		return pistol(human.Pistoling(opponents))
	}

	fmt.Println("Choose two dice from options")
	for j := 1; j < len(rolledDice); j++ {
		fmt.Printf("%d. %s\n", j, rolledDice[j])
	}
	die1, die2 = human.ChooseDiceByIndex()

	if spA == "S" {
		i, t := human.Shoveling(rolledDice[die1], rolledDice[die2])
		shovel(rolledDice, i, t)
	}

	terrain = human.ChooseTerrain(rolledDice[die1][2:], rolledDice[die2][2:])

	fmt.Println("Choose Player that you want to interrogate by number")
	for i, opponent := range opponents {
		fmt.Printf("%d. Player%s\n", i+1, opponent.No)
	}
	n = human.ChoosePlayerByIndex(len(opponents))

	plr = opponents[n-1].No

	var temp string = "05:" + rolledDice[die1] + "," + rolledDice[die2] + "," + terrain + ",P" + plr + "," + spA
	//var temp string = "05:" + "NNB" + "," + "NNB" + "," + "B" + ",P" + plr + "," + spA

	return temp
}

func interrogationReport(args string) string {
	stringSlice := strings.Split(args, ":")
	stringSlice2 := strings.Split(stringSlice[1], ",")

	if stringSlice2[4][1:] != p.No {
		tkncmp.PlayerReportCompute(stringSlice2, &p, opponents)
	}

	// When using shovel
	if stringSlice[0] == "15" {
		for i, _ := range opponents {
			if opponents[i].No == stringSlice2[5][1:] {
				opponents[i].UseAbility("S")
			}
		}
		fmt.Printf("Player %s used shovel.\n", stringSlice2[5][1:])
	}
	fmt.Printf("%s asks %s how many locations they've searched between %s and %s in %s terrain.\n",
		stringSlice2[5], stringSlice2[4], stringSlice2[0][:2], stringSlice2[1][:2], stringSlice2[2])

	fmt.Printf("%s responds %s.\n",
		stringSlice2[4], stringSlice2[3])

	return ""
}

func pistolReport(args string) string {
	message := strings.Split(args[3:], ",")

	for i, _ := range opponents {
		if opponents[i].No == message[0][1:] {
			opponents[i].UseAbility("P")
		}
	}
	fmt.Printf("Player %s used pistol to Player %s\n", message[0][1:], message[1][1:])

	return ""
}

func guessTokens(a1 string, a2 string) string {
	fmt.Println("Submitting a guess:", a1, a2)
	var temp string = "07:P" + p.No + "," + a1 + "," + a2

	return temp
}

func guessCorrect(args string) string {
	stringSlice := strings.Split(args, ":")
	stringSlice2 := strings.Split(stringSlice[1], ",")

	fmt.Printf("Player %s is correct! They have won the game.\n",
		stringSlice2[0])
	fmt.Printf("The treasures were located at %s and %s.\n",
		stringSlice2[1], stringSlice2[2])
	os.Exit(0)

	return ""
}

func guessIncorrect(args string) string {
	message := strings.Split(args, ":")

	fmt.Printf("Player %s is submitting a guess at the treasure locations! Player %s was wrong. They are now disqualified from winning.\n",
		message[1], message[1])

	return ""
}

// Reads from "toPN" named pipe
func readFromPipe(fd *os.File, rd *bufio.Reader) string {
	buff, err := rd.ReadString('\n')
	if err == io.EOF {
		return "exit"
	}
	if err != nil {
		fmt.Println(err)
	}
	if len(buff) > 0 {
		return buff
	}

	return ""
}

// Writes to "fromPN" named pipe
func writeToPipe(fd1 *os.File, args string) {
	fd1.Write([]byte(args))
}

func randInt(min int, max int) int {
	// return min + rand.Intn(max-min)
	return rand.Intn(max-min+1) + min
}

func Find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

func main() {
	var pNo, pipeName, toPN, fromPN string

	pNo = "2"
	pipeName = "all"
	g.autopilot = true
	flag.StringVar(&pNo, "n", "2", "player number")
	flag.StringVar(&pipeName, "pn", "all", "pipe name")
	flag.BoolVar(&g.autopilot, "a", true, "autopilot")
	flag.Parse()
	toPN = "/tmp/" + pipeName + "toP" + pNo
	fromPN = "/tmp/" + pipeName + "fromP" + pNo
	fmt.Println(toPN, fromPN)
	p = player.NewPlayer(pNo)

	// Control logic
	fd, err := os.OpenFile(toPN, os.O_RDONLY, os.ModeNamedPipe) // opens toPN named pipe
	if err != nil {
		fmt.Println(err)
	}
	rd := bufio.NewReader(fd)

	fd1, err1 := os.OpenFile(fromPN, os.O_WRONLY, 0) // opens fromPN named pipe
	if err1 != nil {
		fmt.Println(err1)
	}

	for {
		serverSaid := readFromPipe(fd, rd)
		if serverSaid == "exit" {
			break
		}
		playerReply := selectedFunction(functions[serverSaid[:2]], strings.TrimSpace(serverSaid))
		if playerReply != "" {
			writeToPipe(fd1, playerReply)
		}
	}
}
