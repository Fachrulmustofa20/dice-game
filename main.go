package main

import (
	"fmt"
	"math/rand"
)

type Player struct {
	ID    int
	Dice  []int
	Score int
}

var playerUpdate map[int][]int

// Game function running dice game
func Game(numPlayers, numDice int) {
	players := initializePlayers(numPlayers, numDice)

	for round := 1; ; round++ {
		playerUpdate = make(map[int][]int)
		fmt.Printf("==================\nGiliran %d lempar dadu:\n", round)
		performRoll(players)

		fmt.Println("Setelah evaluasi:")
		activePlayers := evaluateRound(players)
		displayPlayers(activePlayers)

		if len(activePlayers) == 1 {
			fmt.Printf("==================\nGame berakhir karena hanya pemain #%d yang memiliki dadu.\n", activePlayers[0].ID)
			fmt.Printf("Game dimenangkan oleh pemain #%d karena memiliki poin lebih banyak dari pemain lainnya.\n", findWinner(activePlayers))
			break
		}
	}
}

func initializePlayers(numPlayers, numDice int) []*Player {
	players := make([]*Player, numPlayers)

	for i := 0; i < numPlayers; i++ {
		dice := make([]int, numDice)
		players[i] = &Player{
			ID:   i + 1,
			Dice: dice,
		}
	}

	return players
}

// function generate dice every player
func performRoll(players []*Player) {
	for _, player := range players {
		for i := range player.Dice {
			player.Dice[i] = rand.Intn(6) + 1
		}
		fmt.Printf("Pemain #%d (%d): %v\n", player.ID, player.Score, player.Dice)
	}
}

// function evaluate result from dice round
func evaluateRound(players []*Player) []*Player {
	activePlayers := make([]*Player, 0)

	for i, player := range players {
		if len(player.Dice) > 0 {
			player, activePlayers = evaluatePlayer(player, players, activePlayers, i)
		}
	}

	var playersResult = updatePlayerDice(activePlayers)

	return playersResult
}

// update player dice based on playerUpdate
func updatePlayerDice(players []*Player) []*Player {
	for i, player := range players {
		if updateValue, found := playerUpdate[player.ID]; found {
			// Directly update the player's Dice with updateValue
			players[i].Dice = append(players[i].Dice, updateValue...)
		}
	}

	return players
}

// evaluate result from dice every player
func evaluatePlayer(player *Player, allPlayers, activePlayers []*Player, index int) (*Player, []*Player) {
	newActivePlayers := make([]*Player, len(activePlayers))
	copy(newActivePlayers, activePlayers)

	for _, otherPlayer := range allPlayers {
		if otherPlayer.ID != player.ID && len(otherPlayer.Dice) > 0 {
			for j := len(player.Dice) - 1; j >= 0; j-- {
				die := player.Dice[j]
				switch die {
				case 6:
					player.Score++
					player.Dice = removeDie(player.Dice, j)
				case 1:
					var nextPlayerIndex int
					if (index + 1) == len(allPlayers) {
						nextPlayerIndex = 1
					} else {
						nextPlayerIndex = (index + 2)
					}
					playerUpdate[nextPlayerIndex] = append(playerUpdate[nextPlayerIndex], 1)
					player.Dice = removeDie(player.Dice, j)
				}
			}
		}
	}

	if len(player.Dice) > 0 {
		newActivePlayers = append(newActivePlayers, player)
	}

	return player, newActivePlayers
}

// display information player status
func displayPlayers(players []*Player) {
	for _, player := range players {
		fmt.Printf("Pemain #%d (%d): %v\n", player.ID, player.Score, player.Dice)
	}
}

// function to find player winner with the highest point
func findWinner(players []*Player) int {
	winnerID := 0
	maxScore := 0

	for _, player := range players {
		if player.Score > maxScore {
			maxScore = player.Score
			winnerID = player.ID
		}
	}

	return winnerID
}

// function remove slice from dice
func removeDie(dice []int, index int) []int {
	return append(dice[:index], dice[index+1:]...)
}

func main() {
	numPlayers := 3
	numDice := 4
	Game(numPlayers, numDice)
}
