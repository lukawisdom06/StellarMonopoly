package main

import (
	"fmt"
	"github.com/lukawisdom06/stellarMonopoly/game"
)

func main() {
	fmt.Println("🎲 Welcome to Stellar Monopoly! 🎲")
	fmt.Println("================================\n")

	// Initialize the game
	g := game.NewGame()

	// Display game info
	fmt.Printf("Game initialized with %d players\n", len(g.Players))
	fmt.Printf("Board has %d spaces\n\n", len(g.Board.Spaces))

	// Start the game
	g.Play()
}
