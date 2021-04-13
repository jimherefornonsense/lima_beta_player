package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"lima_beta_player/computer/tkncmp"
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

type Combinations struct{
	No         string
	regions    [2]int
	num_tokens int
	potentials map[string][]string
}

type map_entry  struct {
	val int
	key string
}

type map_entries []map_entry

func (s map_entries) Len() int { return len(s) }
func (s map_entries) Less(i, j int) bool { return s[i].val < s[j].val }
func (s map_entries) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

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
		var response string
		p.DisplayTable()
		for _, opponent := range opponents {
			opponent.DisplayTable()
		}
		fmt.Println("Would you like to guess? Y/N")
		fmt.Scanln(&response)
		response = strings.ToUpper(response)
		for response != "Y" && response != "N" {
			fmt.Println("Would you like to guess? Y/N")
			fmt.Scanln(&response)
			response = strings.ToUpper(response)
		}
		if response == "Y" {
			return guessTokens()
		} else {
			fmt.Println("Choose any two dice options from the following or choose A")
			for j := 1; j < len(message); j++ {
				fmt.Printf("%d. %s\n", j, message[j])
			}
		}
	}
	return chooseDice(args)
}

func terrainParser(t1 string, t2 string) string {
	terrainMap := map[string]string{"B": "Beach", "F": "Forest", "M": "Mountain", "A": "All terrians"}
	var t string = "A"
	if g.autopilot {
		var terrain_keys = [4]string{"B", "F", "M", "A"}
		if t1 == "W" && t2 == "W" {
			fmt.Println("first:")
			for k, v := range terrainMap {
				fmt.Printf("%s: %s\n", k, v)
			}
			randomIndex := rand.Intn(len(terrain_keys))
			t = terrain_keys[randomIndex]
		} else if t1 == "W" && t2 != "W" {
			fmt.Println("second:")
			randomIndex := rand.Intn(len(terrain_keys))
			t = terrain_keys[randomIndex]
			for t != "A" && t != t2 {
				randomIndex := rand.Intn(len(terrain_keys))
				t = terrain_keys[randomIndex]
			}
		} else if t1 != "W" && t2 == "W" {
			fmt.Println("third:")
			randomIndex := rand.Intn(len(terrain_keys))
			t = terrain_keys[randomIndex]

			for t != "A" && t != t1 {
				randomIndex := rand.Intn(len(terrain_keys))
				t = terrain_keys[randomIndex]
			}
		} else if t1 == t2 {
			fmt.Println("fourth:")
			return t1
		}
		return t
	} else {
		if t1 == "W" && t2 == "W" {
			fmt.Println("Choose Terrian:")
			for k, v := range terrainMap {
				fmt.Printf("%s: %s\n", k, v)
			}
			fmt.Scanf("%s", &t)
			t = strings.ToUpper(t)
			_, found := terrainMap[t]
			for !found {
				fmt.Println("Invalid terrain, Choose Terrian:")
				fmt.Scanf("%s", &t)
				t = strings.ToUpper(t)
				_, found = terrainMap[t]
			}
		} else if t1 == "W" && t2 != "W" {
			fmt.Println("Choose Terrian:")
			fmt.Printf("%s: %s\n", t2, terrainMap[t2])
			fmt.Printf("A: %s\n", terrainMap["A"])
			fmt.Scanf("%s", &t)
			t = strings.ToUpper(t)
			for t != "A" && t != t2 {
				fmt.Println("Invalid terrain, Choose Terrian:")
				fmt.Scanf("%s", &t)
				t = strings.ToUpper(t)
			}
		} else if t1 != "W" && t2 == "W" {
			fmt.Println("Choose Terrian:")
			fmt.Printf("%s: %s\n", t1, terrainMap[t1])
			fmt.Printf("A: %s\n", terrainMap["A"])
			fmt.Scanf("%s", &t)
			t = strings.ToUpper(t)
			for t != "A" && t != t1 {
				fmt.Println("Invalid terrain, Choose Terrian:")
				fmt.Scanf("%s", &t)
				t = strings.ToUpper(t)
			}
		} else if t1 == t2 {
			return t1
		}
		return t
	}
}

func terrainParser_comp(t1 string, t2 string, cm Combinations) string {
	var t string = "A"
	var m = make(map[string]int)
 m["A"] =  len(cm.potentials["A"])	
 m["B"] =  len(cm.potentials["B"])
 m["F"] =  len(cm.potentials["F"])
 m["M"] =  len(cm.potentials["M"])
	
	

	var es map_entries
	for k, v := range m {
		es = append(es, map_entry{val: v, key: k})
	}
	
	sort.Sort(sort.Reverse(es))
	

	if t1 == "W" && t2 == "W" {
		
		t = es[0].key
	//t = "A"
	return t
	} else if t1 == "W" && t2 != "W" {
		
		t = "A"
	} else if t1 != "W" && t2 == "W" {
		fmt.Println("third:")
		t = "A"
	} else if t1 == t2 {
		fmt.Println("fourth:")
		return t1
	}
	
//	return  es[0].key
	return t
}

func chooseDice(args string) string {
	rolledDice := strings.Split(args[3:], ",")
	var n int
	var die1 int = 0
	var die2 int = 0
	var terrain string

	// Temporary using random choosing, still implementing intelligent one
	if g.autopilot {
		var plr string
		message := strings.Split(args[3:], ",")
		

		var combinations_group []Combinations

		opponents_sets := make(map[string][]string)

		
		for j := range opponents {
			for i := 1; i <= 3; i++  {
				for k := i+1; k <= 3; k++  {
			var c1 Combinations 
			c1.potentials = make(map[string][]string)
			c1.potentials["A"] = opponents[j].UnfirmedOneTokensInRegion(message[i][:2], message[k][:2], "A")
			c1.potentials["B"] = opponents[j].UnfirmedOneTokensInRegion(message[i][:2], message[k][:2], "B")
			c1.potentials["F"] = opponents[j].UnfirmedOneTokensInRegion(message[i][:2], message[k][:2], "F")
			c1.potentials["M"] = opponents[j].UnfirmedOneTokensInRegion(message[i][:2], message[k][:2], "M")

			c1.num_tokens = tkncmp.NumTknsInRegion(message[i][:2], message[k][:2], "A")
			c1.regions = [2]int{i, k}
			c1.No = opponents[j].No

			var c2 Combinations 
			c2.potentials = make(map[string][]string)
			c2.potentials["A"] = opponents[j].UnfirmedOneTokensInRegion(message[k][:2], message[i][:2], "A")
			c2.potentials["B"] = opponents[j].UnfirmedOneTokensInRegion(message[k][:2], message[i][:2], "B")
			c2.potentials["F"] = opponents[j].UnfirmedOneTokensInRegion(message[k][:2], message[i][:2], "F")
			c2.potentials["M"] = opponents[j].UnfirmedOneTokensInRegion(message[k][:2], message[i][:2], "M")

			c2.num_tokens = tkncmp.NumTknsInRegion(message[k][:2], message[i][:2], "A")
			c2.regions = [2]int{k, i}
			c2.No = opponents[j].No
		
			combinations_group = append(combinations_group,c1)
			combinations_group = append(combinations_group,c2)
				}
		
		}
		var temp []string
		for l := 1; l <= len(opponents[j].PotentialObtainedTknsList); l++ {
	
			for _, set := range opponents[j].PotentialObtainedTknsList[l] {
				setString :=  fmt.Sprintf("%s", set)
				temp = append(temp,setString)
			}
		} 
		opponents_sets[opponents[j].No] = temp
		}

		sort.SliceStable(combinations_group, func(i, j int) bool {
			
				return  len(combinations_group[i].potentials["A"]) >  len(combinations_group[j].potentials["A"])

			
			
//			return  len(combinations_group[i].potentials["A"]) > 0 && len(combinations_group[j].potentials["A"]) > 0 && combinations_group[i].num_tokens > combinations_group[j].num_tokens
		})
		
         var decided_player Combinations
	for i :=0; i<len(combinations_group); i++ {
		_, found := Find(opponents_sets[combinations_group[i].No], fmt.Sprintf("%s", combinations_group[i].potentials["A"]))
		if !found {
			fmt.Println("check")
			plr = combinations_group[i].No
			die1 = combinations_group[i].regions[0]
			die2 = combinations_group[i].regions[1]
			decided_player = combinations_group[i]
			break
		} else{
			//fmt.Println("found duplicate")
		//	fmt.Println("first")
		//	fmt.Println(opponents_sets[combinations_group[i].No])
		//	fmt.Println("second")
		//	setString := fmt.Sprintf("%s", combinations_group[i].potentials)
		//	for j :=0; j<len(combinations_group[i].potentials["A"]); j++ {
		//	fmt.Println(combinations_group[i].potentials["A"][j])
		//	}
			//fmt.Println()
		}

	}
	


if die1==0 || die2==0 {
	sort.SliceStable(combinations_group, func(i, j int) bool {
			
	//	return  combinations_group[i].num_tokens <  combinations_group[i].num_tokens

	return  len(combinations_group[i].potentials["A"]) >  len(combinations_group[j].potentials["A"])

	
//			return  len(combinations_group[i].potentials["A"]) > 0 && len(combinations_group[j].potentials["A"]) > 0 && combinations_group[i].num_tokens > combinations_group[j].num_tokens
})
plr = combinations_group[0].No
			die1 = combinations_group[0].regions[0]
			die2 = combinations_group[0].regions[1]
			decided_player = combinations_group[0]
	fmt.Println("heyyy")
}
	fmt.Println(die1)
	fmt.Println(die2)
	fmt.Println(plr)
		
		 
	
		terrain = terrainParser_comp(string(rolledDice[die1][2]), string(rolledDice[die2][2]),decided_player)
	 /*   
		if opponents_intials[0] < opponents_intials[1] {
		plr = opponents[0].No
		} else {
			plr = opponents[1].No
		}
		*/
	/*	plr = randInt(1, len(opponents)+1)
		for strconv.Itoa(plr) == p.No {
			plr = randInt(1, len(opponents)+1)
		}*/
		var temp string = "05:" + rolledDice[die1] + "," + rolledDice[die2] + "," + terrain + ",P" + plr
		
//	fmt.Println(strconv.Itoa(plr))
		fmt.Println(len(opponents))
		fmt.Println("temp-" + temp)
		return temp
	} else {
		var plr string
		fmt.Println("Choose first die by number")
		_, err := fmt.Scan(&n)
		for err != nil || n > 3 || n < 1 {
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("Out of range, choose first die by number")
			}
			_, err = fmt.Scan(&n)
		}
		die1 = n
		//fmt.Println("die1:", die1)

		fmt.Println("Choose second die by number")
		_, err = fmt.Scan(&n)
		for err != nil || n > 3 || n < 1 || n == die1 {
			if err != nil {
				fmt.Println(err)
			} else if n > 3 || n < 1 {
				fmt.Println("Out of range, choose second die by number")
			} else if n == die1 {
				fmt.Println("Die has chosen, enter another number for second die")
			}
			_, err = fmt.Scan(&n)
		}
		die2 = n
		//fmt.Println("die2:", die2)

		terrain = terrainParser(rolledDice[die1][2:], rolledDice[die2][2:])

		fmt.Println("Choose Player that you want to interrogate by number")
		for i, opponent := range opponents {
			fmt.Printf("%d. Player%s\n", i+1, opponent.No)
		}
		_, err = fmt.Scan(&n)
		for err != nil || n > len(opponents) || n < 1 {
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("Out of range, choose Player that you want to interrogate by number")
			}
			_, err = fmt.Scan(&n)
		}
		plr = opponents[n-1].No

		var temp string = "05:" + rolledDice[die1] + "," + rolledDice[die2] + "," + terrain + ",P" + plr
		return temp
	}
}

func interrogationReport(args string) string {
	stringSlice := strings.Split(args, ":")
	stringSlice2 := strings.Split(stringSlice[1], ",")

	if stringSlice2[4][1:] != p.No {
		tkncmp.PlayerReportCompute(stringSlice2, &p, opponents)
	}

	fmt.Printf("%s asks %s how many locations they've searched between %s and %s in %s terrain.\n",
		stringSlice2[5], stringSlice2[4], stringSlice2[0][:2], stringSlice2[1][:2], stringSlice2[2])

	fmt.Printf("%s responds %s.\n",
		stringSlice2[4], stringSlice2[3])
	return ""
}

func isValidToken(token string) bool {
	for _, t := range tkncmp.TokenMap {
		if strings.ToUpper(token) == t {
			return true
		}
	}
	return false
}

func guessTokens(answer ...string) string {
	if g.autopilot {
		fmt.Println("Submitting a guess:", answer[0], answer[1])
		var temp string = "07:P" + p.No + "," + answer[0] + "," + answer[1]
		return temp
	} else {
		var first_token string
		var second_token string

		fmt.Println("Choose first token: ")
		fmt.Scanf("%s", &first_token)
		for !isValidToken(first_token) {
			fmt.Println("Invalid token, please choose the first token: ")
			fmt.Scanf("%s", &first_token)
		}
		fmt.Println("Choose second token: ")
		fmt.Scanf("%s", &second_token)
		for !isValidToken(second_token) || second_token == first_token {
			if second_token == first_token {
				fmt.Println("Same guess, please choose the second token: ")
			} else {
				fmt.Println("Invalid token, please choose the second token: ")
			}
			fmt.Scanf("%s", &second_token)
		}
		var temp string = "07:P" + p.No + "," + strings.ToUpper(first_token) + "," + strings.ToUpper(second_token)
		return temp
	}
}

func guessCorrect(args string) string {
	stringSlice := strings.Split(args, ":")
	stringSlice2 := strings.Split(stringSlice[1], ",")
	fmt.Printf("Player %s is correct! They have won the game.\n",
		stringSlice2[0])
	fmt.Printf("The treasures were located at %s and %s.\n",
		stringSlice2[1], stringSlice2[2])
	return ""
}

func guessIncorrect(args string) string {
	message := strings.Split(args, ":")
	fmt.Printf("Player %s is submitting a guess at the treasure locations!. Player %s was wrong. They are now disqualified from winning.\n",
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
