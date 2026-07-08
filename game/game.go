package game

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Game represents the main Monopoly game
type Game struct {
	Players      []*Player
	Board        *Board
	CurrentTurn  int
	GameOver     bool
	RoundCount   int
	PlayerCount  int
}

// NewGame initializes a new Monopoly game
func NewGame() *Game {
	board := NewBoard()
	
	// Get number of players
	playerCount := getPlayerCount()
	
	players := make([]*Player, playerCount)
	for i := 0; i < playerCount; i++ {
		players[i] = NewPlayer(i+1, fmt.Sprintf("Player %d", i+1), 1500)
	}

	return &Game{
		Players:     players,
		Board:       board,
		CurrentTurn: 0,
		GameOver:    false,
		RoundCount:  0,
		PlayerCount: playerCount,
	}
}

// Play starts the main game loop
func (g *Game) Play() {
	fmt.Println("Game Start! Each player starts with $1500")
	fmt.Println("========================================\n")

	maxRounds := 50 // Prevent infinite games during testing

	for !g.GameOver && g.RoundCount < maxRounds {
		g.RoundCount++
		fmt.Printf("\n--- Round %d ---\n", g.RoundCount)

		for i := 0; i < len(g.Players); i++ {
			if g.Players[i].IsBankrupt {
				continue
			}

			currentPlayer := g.Players[i]
			fmt.Printf("\n%s's turn (Balance: $%d)\n", currentPlayer.Name, currentPlayer.Balance)
			fmt.Println("1. Roll Dice")
			fmt.Println("2. View Properties")
			fmt.Println("3. View Balance")
			fmt.Println("4. End Turn")

			choice := getUserInput("Choose action: ")

			switch choice {
			case "1":
				g.rollDiceAndMove(currentPlayer)
			case "2":
				currentPlayer.DisplayProperties()
			case "3":
				fmt.Printf("Current Balance: $%d\n", currentPlayer.Balance)
			case "4":
				continue
			default:
				fmt.Println("Invalid choice. Ending turn.")
			}

			// Check if only one player remains
			activePlayers := g.getActivePlayers()
			if len(activePlayers) == 1 {
				g.GameOver = true
				break
			}
		}
	}

	g.EndGame()
}

// rollDiceAndMove handles dice rolling and movement
func (g *Game) rollDiceAndMove(player *Player) {
	dice := RollDice()
	fmt.Printf("%s rolled: %d + %d = %d\n", player.Name, dice.Die1, dice.Die2, dice.Total)

	// Move player
	newPosition := (player.Position + dice.Total) % len(g.Board.Spaces)
	
	// Check if player passed GO
	if newPosition < player.Position {
		fmt.Printf("%s passed GO! Collecting $200\n", player.Name)
		player.Balance += 200
	}

	player.Position = newPosition
	space := g.Board.Spaces[newPosition]

	fmt.Printf("%s moved to: %s\n", player.Name, space.Name)

	// Handle space action
	g.handleSpaceAction(player, space)
}

// handleSpaceAction determines what happens when a player lands on a space
func (g *Game) handleSpaceAction(player *Player, space *Space) {
	switch space.Type {
	case "Property":
		g.handlePropertyLanding(player, space)
	case "Railroad":
		g.handleRailroadLanding(player, space)
	case "Utility":
		g.handleUtilityLanding(player, space)
	case "Tax":
		fmt.Printf("%s pays $%d in taxes\n", player.Name, space.Cost)
		player.Balance -= space.Cost
	case "GO":
		fmt.Printf("%s is at GO! No action needed.\n", player.Name)
	case "FreeParking":
		fmt.Printf("%s is at Free Parking.\n", player.Name)
	case "GoToJail":
		fmt.Printf("%s has been sent to Jail!\n", player.Name)
		player.InJail = true
		player.Position = 10
	case "Jail":
		fmt.Printf("%s is in Jail.\n", player.Name)
	}

	// Check if player is bankrupt
	if player.Balance < 0 {
		fmt.Printf("%s is BANKRUPT!\n", player.Name)
		player.IsBankrupt = true
	}
}

// handlePropertyLanding handles landing on a property
func (g *Game) handlePropertyLanding(player *Player, space *Space) {
	if space.Owner == nil {
		fmt.Printf("%s landed on unowned property: %s (Cost: $%d)\n", player.Name, space.Name, space.Cost)
		buyChoice := getUserInput("Do you want to buy this property? (yes/no): ")
		if strings.ToLower(buyChoice) == "yes" && player.Balance >= space.Cost {
			player.Balance -= space.Cost
			space.Owner = player
			player.Properties = append(player.Properties, space)
			fmt.Printf("%s bought %s for $%d\n", player.Name, space.Name, space.Cost)
		}
	} else if space.Owner.Name != player.Name {
		rent := space.Rent
		fmt.Printf("%s must pay $%d rent to %s\n", player.Name, rent, space.Owner.Name)
		player.Balance -= rent
		space.Owner.Balance += rent
	}
}

// handleRailroadLanding handles landing on a railroad
func (g *Game) handleRailroadLanding(player *Player, space *Space) {
	if space.Owner == nil {
		fmt.Printf("%s landed on unowned railroad: %s (Cost: $%d)\n", player.Name, space.Name, space.Cost)
		buyChoice := getUserInput("Do you want to buy this railroad? (yes/no): ")
		if strings.ToLower(buyChoice) == "yes" && player.Balance >= space.Cost {
			player.Balance -= space.Cost
			space.Owner = player
			player.Properties = append(player.Properties, space)
			fmt.Printf("%s bought %s for $%d\n", player.Name, space.Name, space.Cost)
		}
	} else if space.Owner.Name != player.Name {
		rent := 25 * len(space.Owner.GetRailroads())
		fmt.Printf("%s must pay $%d rent to %s\n", player.Name, rent, space.Owner.Name)
		player.Balance -= rent
		space.Owner.Balance += rent
	}
}

// handleUtilityLanding handles landing on a utility
func (g *Game) handleUtilityLanding(player *Player, space *Space) {
	if space.Owner == nil {
		fmt.Printf("%s landed on unowned utility: %s (Cost: $%d)\n", player.Name, space.Name, space.Cost)
		buyChoice := getUserInput("Do you want to buy this utility? (yes/no): ")
		if strings.ToLower(buyChoice) == "yes" && player.Balance >= space.Cost {
			player.Balance -= space.Cost
			space.Owner = player
			player.Properties = append(player.Properties, space)
			fmt.Printf("%s bought %s for $%d\n", player.Name, space.Name, space.Cost)
		}
	} else if space.Owner.Name != player.Name {
		dice := RollDice()
		rent := dice.Total * 4
		fmt.Printf("%s rolled %d. Must pay $%d rent to %s\n", player.Name, dice.Total, rent, space.Owner.Name)
		player.Balance -= rent
		space.Owner.Balance += rent
	}
}

// getActivePlayers returns list of non-bankrupt players
func (g *Game) getActivePlayers() []*Player {
	var active []*Player
	for _, p := range g.Players {
		if !p.IsBankrupt {
			active = append(active, p)
		}
	}
	return active
}

// EndGame handles game conclusion
func (g *Game) EndGame() {
	fmt.Println("\n\n========== GAME OVER ==========")
	
	var winner *Player
	maxBalance := -1
	
	for _, player := range g.Players {
		if player.Balance > maxBalance {
			maxBalance = player.Balance
			winner = player
		}
	}

	if winner != nil {
		fmt.Printf("🎉 %s WINS with $%d! 🎉\n", winner.Name, winner.Balance)
	}

	fmt.Println("\nFinal Standings:")
	for i, player := range g.Players {
		status := "Active"
		if player.IsBankrupt {
			status = "Bankrupt"
		}
		fmt.Printf("%d. %s - $%d (%s)\n", i+1, player.Name, player.Balance, status)
	}
}

// Helper functions
func getPlayerCount() int {
	for {
		fmt.Print("Enter number of players (2-4): ")
		input := getUserInput("")
		count, err := strconv.Atoi(input)
		if err == nil && count >= 2 && count <= 4 {
			return count
		}
		fmt.Println("Invalid input. Please enter 2, 3, or 4.")
	}
}

func getUserInput(prompt string) string {
	if prompt != "" {
		fmt.Print(prompt)
	}
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}
