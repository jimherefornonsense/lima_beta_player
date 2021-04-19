package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"lima_beta_player/computer/selectorcmp"
	"lima_beta_player/computer/tkncmp"
	"lima_beta_player/human"
	"lima_beta_player/player"
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
	"13": barrelReport,
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

	fmt.Println("Using pistol!!!!!!!!!!!!!")
	msg = "05:" + d1 + "," + d2 + "," + t + "," + "P" + pNo + "," + "P"

	return msg
}

func shovel(selectedDice []string, index int, swappedTerrain string) {
	fmt.Println("Using shovel!!!!!!!!!!!!!")
	selectedDice[index] = strings.Replace(selectedDice[index], selectedDice[index][2:], swappedTerrain, 1)
}

func barrel(rerolledDice []int) string {
	var msg string

	fmt.Println("Using barrel!!!!!!!!!!!!!")
	sort.Ints(rerolledDice)
	var suffix string = ""
	for k := 0; k < len(rerolledDice); k++ {
		suffix += ","
		suffix += strconv.Itoa(rerolledDice[k])
	}

	msg = "12:" + strconv.Itoa(len(rerolledDice)) + suffix

	return msg
}

func chooseDice(args string) string {
	rolledDice := strings.Split(args[3:], ",")
	var n int
	var die1 int = 0
	var die2 int = 0
	var terrain, plr string

	if g.autopilot {
		var comp_SPA string = "Q"
		var aCombination selectorcmp.Combination
		var isIntelligent bool

		isIntelligent, aCombination = selectorcmp.Selection(rolledDice, opponents)
		if isIntelligent {
			doIt, targetNo := selectorcmp.IsUsePistol(p, opponents)
			if doIt {
				if p.UseAbility("P") {
					return pistol(selectorcmp.Pistoling(rolledDice, targetNo, opponents))
				}
			}
		} else {
			if p.UseAbility("B") {
				return barrel(selectorcmp.Barreling(rolledDice, aCombination))
			}
		}

		plr = aCombination.No
		die1 = aCombination.Region[0]
		die2 = aCombination.Region[1]
		terrain = aCombination.Terrain

		if rolledDice[die1][2:] == rolledDice[die2][2:] && rolledDice[die1][2:] != "W" {
			if p.UseAbility("S") {
				shovel(rolledDice, die1, selectorcmp.Shoveling(rolledDice[die1]))
				_, aCombination = selectorcmp.Selection(rolledDice, opponents)
				comp_SPA = "S"
				plr = aCombination.No
				die1 = aCombination.Region[0]
				die2 = aCombination.Region[1]
				terrain = aCombination.Terrain
			}
		}

		var temp string = "05:" + rolledDice[die1] + "," + rolledDice[die2] + "," + terrain + ",P" + plr + "," + comp_SPA
		return temp
	}

	// Manual
	spA := human.IsUsingSpA(&p)

	if spA == "P" {
		return pistol(human.Pistoling(opponents))
	} else if spA == "B" {
		return barrel(human.Barreling(rolledDice[1], rolledDice[2], rolledDice[3]))
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

func barrelReport(args string) string {
	message := strings.Split(args[3:], ",")
	for i, _ := range opponents {
		if opponents[i].No == message[0][1:] {
			opponents[i].UseAbility("B")
		}
	}
	fmt.Println("Player " + message[0][1:] + " has rerolled " + message[1] + "," + message[2] + "," + message[3])

	return chooseDice(args)
}

func shovelReport(args string) string {
	message := strings.Split(args[3:], ",")
	for i, _ := range opponents {
		if opponents[i].No == message[0][1:] {
			opponents[i].UseAbility("S")
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
	fmt.Println("written msg")
	fmt.Println(args)
	fd1.Write([]byte(args))
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
			fmt.Println("PlayerReply")
			fmt.Println(playerReply)
			writeToPipe(fd1, playerReply)
		}
	}
}
