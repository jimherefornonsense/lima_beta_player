package human

import (
	"fmt"
	"lima_beta_player/player"
	"strings"
)

func IsGuessing() bool {
	var response string
	fmt.Println("Would you like to guess? Y/N")
	fmt.Scanln(&response)
	response = strings.ToUpper(response)
	for response != "Y" && response != "N" {
		fmt.Println("Would you like to guess? Y/N")
		fmt.Scanln(&response)
		response = strings.ToUpper(response)
	}

	if response == "Y" {
		return true
	}

	return false
}

func isValidToken(token string) bool {
	tokenSet := "1B1F1M2B2F2M3B3F3M4B4F4M5B5F5M6B6F6M7B7F7M8B8F8M"
	var found bool = false

	if strings.TrimSpace(token) != "" {
		found = strings.Contains(tokenSet, token)
	}

	return found
}

func GuessTreasures() (string, string) {
	var firstToken string
	var secondToken string

	fmt.Println("Choose first token: ")
	fmt.Scanf("%s", &firstToken)
	firstToken = strings.ToUpper(firstToken)
	for !isValidToken(firstToken) {
		fmt.Println("Invalid token, please choose the first token: ")
		fmt.Scanf("%s", &firstToken)
		firstToken = strings.ToUpper(firstToken)
	}
	fmt.Println("Choose second token: ")
	fmt.Scanf("%s", &secondToken)
	secondToken = strings.ToUpper(secondToken)
	for !isValidToken(secondToken) || secondToken == firstToken {
		if secondToken == firstToken {
			fmt.Println("Same guess, please choose the second token: ")
		} else {
			fmt.Println("Invalid token, please choose the second token: ")
		}
		fmt.Scanf("%s", &secondToken)
		secondToken = strings.ToUpper(secondToken)
	}

	return firstToken, secondToken
}

func ChooseTerrain(t1 string, t2 string) string {
	terrainMap := map[string]string{"B": "Beach", "F": "Forest", "M": "Mountain", "A": "All terrians"}
	var t string = "A"

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

func ChooseDiceByIndex() (int, int) {
	var n, die1, die2 int

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

	return die1, die2
}

func ChoosePlayerByIndex(numPlayers int) int {
	var n int

	_, err := fmt.Scan(&n)
	for err != nil || n > numPlayers || n < 1 {
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Out of range, choose a Player that you want to interrogate to by number")
		}
		_, err = fmt.Scan(&n)
	}

	return n
}

func IsUsingSpA(plr *player.Player) string {
	var a string
	var optionSet string = "Q|"

	fmt.Println("Do you want to use special ability?")
	fmt.Println("Q: no")
	if plr.Pistol != 0 {
		fmt.Println("P: Pistol")
		optionSet = optionSet + "P|"
	}
	if plr.Shovel != 0 {
		fmt.Println("S: Shovel")
		optionSet = optionSet + "S|"
	}
	if plr.Barrel != 0 {
		fmt.Println("B: Barrel")
		optionSet = optionSet + "B|"
	}

	_, err := fmt.Scanln(&a)
	a = strings.ToUpper(a)
	for err != nil || !strings.Contains(optionSet, a) {
		fmt.Println("Choose the right option!")
		_, err = fmt.Scanln(&a)
		a = strings.ToUpper(a)
	}

	if plr.UseAbility(a) {
		return a
	}

	return "Q"
}

func Pistoling(opponents []player.Player) (string, string, string, string) {
	var start, end, t string
	var index int
	var directionSet string = "NN|NE|EE|SE|SS|SW|WW|NW"

	fmt.Println("Choose two directions of start and end (NN NE EE SE SS SW WW NW NN), and a terrain (B: Beach F: Forest M: Mountain A: All terrains), and a player you want to ask by number (separate by space)")
	for i := range opponents {
		fmt.Printf("%d. Player %s\n", i+1, opponents[i].No)
	}
	fmt.Scanf("%s%s%s%d", &start, &end, &t, &index)
	start = strings.ToUpper(start)
	end = strings.ToUpper(end)
	t = strings.ToUpper(t)
	for !strings.Contains(directionSet, start) || !strings.Contains(directionSet, end) || (t != "B" && t != "F" && t != "M" && t != "A") || index > len(opponents) || index < 1 {
		fmt.Println("Submit the right format!")
		fmt.Scanf("%s%s%s%d", &start, &end, &t, &index)
		start = strings.ToUpper(start)
		end = strings.ToUpper(end)
		t = strings.ToUpper(t)
	}

	return start + "W", end + "W", t, opponents[index-1].No
}

func Shoveling(d1 string, d2 string) (int, string) {
	var i int
	var t string

	fmt.Println("Choose one die by number and a terrain code (B: Beach F: Forest M: Mountain) for wanted terrain")
	fmt.Println("1.", d1)
	fmt.Println("2.", d2)
	fmt.Scanf("%d%s", &i, &t)
	t = strings.ToUpper(t)
	for (i != 1 && i != 2) || (t != "B" && t != "F" && t != "M") {
		fmt.Println("Submit the right format!")
		fmt.Scanf("%d%s", &i, &t)
		t = strings.ToUpper(t)
	}

	return i, t
}

func Barreling(d1 string, d2 string, d3 string) []int {
	var rerolledDice []int
	var num int
	var first int = -1
	var second int = -2

	fmt.Println("Choose how many dice you want to reroll (dice numbers 1-3)")
	fmt.Scanf("%d", &num)
	for num < 1 || num > 3 {
		fmt.Println("Submit the right answer!")
		fmt.Scanf("%d", &num)
	}

	if num == 1 {
		fmt.Println("Select the die options by number")
		fmt.Println("1.", d1)
		fmt.Println("2.", d2)
		fmt.Println("3.", d3)
		fmt.Scanf("%d", &first)
		for first < 1 || first > 3 {
			fmt.Println("Submit the right answer!")
			fmt.Scanf("%d", &first)
		}
		rerolledDice = append(rerolledDice, first-1)
	} else if num == 2 {
		fmt.Println("Select the die options by number (separate by space)")
		fmt.Println("1.", d1)
		fmt.Println("2.", d2)
		fmt.Println("3.", d3)
		fmt.Scanf("%d%d", &first, &second)
		for (first < 1 || first > 3) || (second < 1 || second > 3) || first == second {
			fmt.Println("Submit the right answer!")
			fmt.Scanf("%d%d", &first, &second)
		}
		rerolledDice = append(rerolledDice, first-1)
		rerolledDice = append(rerolledDice, second-1)
	} else {
		rerolledDice = append(rerolledDice, 0)
		rerolledDice = append(rerolledDice, 1)
		rerolledDice = append(rerolledDice, 2)
	}

	return rerolledDice
}
