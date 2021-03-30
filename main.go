package main

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"lima_beta_player/computer/tkncmp"
	"lima_beta_player/player"
)

// Game struct
type Game struct {
	totalPlayers   int
	activePlayers  int
	mode           int
	current_round  int
	maximum_rounds int
	leftTokens     []string
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
	"99": errorMsg,
}

func errorMsg(arg string) string {
	fmt.Println(arg)
	return ""
}

func (gm *Game) playerNO(args string) string {
	gm.totalPlayers, _ = strconv.Atoi(args[len(args)-1:])
	gm.activePlayers = gm.totalPlayers
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
	for _, token := range playerTokens {
		p.PlayerTerrains = append(p.PlayerTerrains, token)
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

	if g.mode == 2 {
		p.DisplayTable()
		for _, opponent := range opponents {
			opponent.DisplayTable()
		}
		

		var final_guesses_first int = 0
		var final_guesses_second int = 0
			 for i := range p.Table {
			 if p.Table[i].Beach==-1  {
				final_guesses_first++;
			 } else if p.Table[i].Forest==-1  {
				final_guesses_first++;
			 } else if p.Table[i].Mountain==-1   {
				final_guesses_first++;
			}


			}
			for i := range p.Table {
				if p.Table[i].Beach==3  {
					final_guesses_second++;
				} else if p.Table[i].Forest==3  {
					final_guesses_second++;
				} else if p.Table[i].Mountain==3  {
					final_guesses_second++;
			   }
   
			   }
if final_guesses_second==2 {
	return guessTokens();
} else if final_guesses_first==1 && final_guesses_second==1 {
	return guessTokens();

} else if final_guesses_second==2  {
	return guessTokens();

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
	if g.mode == 2 {
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

func chooseDice(args string) string {
	g.current_round++
	rolledDice := strings.Split(args[3:], ",")
	var n, die1, die2 int
	var terrain string

	if g.mode == 2 {
		var plr int=0
	    message := strings.Split(args[3:], ",")
		var regions []int
		for j := 1; j < len(message); j++ {
			regions = append(regions,player.TableIndexMap[message[j][:2]]) 
		}

		for j := 1; j < len(regions); j++ {
         fmt.Println(regions[j])
		}
		
			
             
				 var total_potentials1 int = 0
				 var total_potentials2 int = 0

				 var intials1 int = 0
				 var intials2 int = 0

				 var ratios1 float64 = 0
				 var ratios2 float64 = 0
				for i := regions[0]; i<=regions[1]; i++  {
					//for j := regions[0]; i<=regions[1]; i++  {
						for j := range opponents {
                   if opponents[j].Table[i].Beach==2||opponents[j].Table[i].Forest==2||opponents[j].Table[i].Mountain==2 {
					total_potentials1++;
				   } else if opponents[j].Table[i].Beach==-1||opponents[j].Table[i].Forest==-1||opponents[j].Table[i].Mountain==-1 {
					intials1++;
				   }
				}
			}
				

				for i := regions[1]; i<=regions[2]; i++  {
					for j := range opponents {
						if opponents[j].Table[i].Beach==2||opponents[j].Table[i].Forest==2||opponents[j].Table[i].Mountain==2 {
							total_potentials2++;
						   } else if opponents[j].Table[i].Beach==-1||opponents[j].Table[i].Forest==-1||opponents[j].Table[i].Mountain==-1 {
							intials2++;
						   }
					}
				}

				ratios1 = float64(intials1)/float64(intials1+total_potentials1)
				ratios2 = float64(intials2)/float64(intials2+total_potentials2)

if ratios1==ratios2 {
if(regions[1]-regions[0]<regions[2]-regions[1]) {
die1=1
die2=2
} else {
	die1=2
	die2=3
}
} else {
	if ratios1<ratios2 {
		die1=1
		die2=2
	} else {
		die1=2
		die2=3
	}

}
                  
		

		

		terrain = terrainParser(string(rolledDice[die1][2]), string(rolledDice[die2][2]))

		plr = randInt(1, len(opponents)+1)
		for strconv.Itoa(plr) == p.No {
			plr = randInt(1, len(opponents)+1)
		}
		var temp string = "05:" + rolledDice[die1] + "," + rolledDice[die2] + "," + terrain + ",P" + strconv.Itoa(plr)
		fmt.Println(die1)
		fmt.Println(die2)
		fmt.Println(p.No)
		fmt.Println(strconv.Itoa(plr))
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

func guessTokens() string {
	var first_token string
	var second_token string

	var tokens []string
	if g.mode == 2 {
			
		for i := range p.Table {
			if p.Table[i].Beach==-1  {
				tokens = append(tokens,strconv.Itoa(i+1)+"B")

			} else if p.Table[i].Forest==-1  {
				tokens = append(tokens,strconv.Itoa(i+1)+"F")
			} else if p.Table[i].Mountain==-1   {
				tokens = append(tokens,strconv.Itoa(i+1)+"M")
		   }


		   }
		   for i := range p.Table {
			   if p.Table[i].Beach==3  {
				tokens = append(tokens,strconv.Itoa(i+1)+"B")
			   } else if p.Table[i].Forest==3  {
				tokens = append(tokens,strconv.Itoa(i+1)+"F")
			   } else if p.Table[i].Mountain==3  {
				tokens = append(tokens,strconv.Itoa(i+1)+"M")
			  }
  
			  }

		   

/*	
	randomIndex := rand.Intn(len(tkncmp.TokenMap))
		first_token = tkncmp.TokenMap[randomIndex]
		_, found := Find(p.PlayerTerrains, first_token)
		for found {
			randomIndex := rand.Intn(len(tkncmp.TokenMap))
			first_token = tkncmp.TokenMap[randomIndex]
			_, found = Find(p.PlayerTerrains, first_token)
		}

		randomIndex2 := rand.Intn(len(tkncmp.TokenMap))
		second_token = tkncmp.TokenMap[randomIndex2]
		_, found2 := Find(p.PlayerTerrains, second_token)
		for found2 || first_token == second_token {
			randomIndex2 := rand.Intn(len(tkncmp.TokenMap))
			second_token = tkncmp.TokenMap[randomIndex2]
			_, found2 = Find(p.PlayerTerrains, second_token)
		}
*/
		fmt.Println("Submitting a guess:", tokens[0], tokens[1])

		var temp string = "07:P" + p.No + "," + tokens[0]+ "," + tokens[1]
		return temp
	} else {
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
	var mode, rounds int

	// Mode of human-controlled or computer-controlled
	fmt.Println("Do you want player to be controlled by human or by computer? Choose 1 for human and 2 for computer")
	fmt.Println("1.Human")
	fmt.Println("2.Computer")
	fmt.Scanf("%d", &mode)
	for mode != 1 && mode != 2 {
		fmt.Println("Invalid selection")
		fmt.Println("Do you want player to be controlled by human or by computer? Choose 1 for human and 2 for computer")
		fmt.Println("1.Human")
		fmt.Println("2.Computer")
		fmt.Scanf("%d", &mode)
	}
	g.mode = mode

	// Set maximum round
	if g.mode == 2 {
		fmt.Println("Choose the maximum number of rounds you want")
		fmt.Scanf("%d", &rounds)
		for rounds < 1 {
			fmt.Println("You should choose the minimum number of rounds to be 1")
			fmt.Scanf("%d", &rounds)
		}
		// for rounds > 5 {
		// 	fmt.Println("You can choose the maximum number of rounds upto 5")
		// 	fmt.Scanf("%d", &rounds)
		// }
		g.maximum_rounds = rounds
		g.current_round = 0
	}

	// Set player number and directory prefixed Name
	fmt.Println("Enter your Player Number and pipe Prefixed Name: (separated by space)")
	fmt.Scanf("%s%s", &pNo, &pipeName)
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
