package game

import "fmt"

// Player represents a player in the game
type Player struct {
	ID        int
	Name      string
	Balance   int
	Position  int
	Properties []*Space
	InJail    bool
	IsBankrupt bool
}

// NewPlayer creates a new player
func NewPlayer(id int, name string, startingBalance int) *Player {
	return &Player{
		ID:        id,
		Name:      name,
		Balance:   startingBalance,
		Position:  0,
		Properties: make([]*Space, 0),
		InJail:    false,
		IsBankrupt: false,
	}
}

// DisplayProperties shows all properties owned by the player
func (p *Player) DisplayProperties() {
	if len(p.Properties) == 0 {
		fmt.Printf("%s owns no properties.\n", p.Name)
		return
	}

	fmt.Printf("\n%s's Properties:\n", p.Name)
	totalValue := 0
	for _, prop := range p.Properties {
		fmt.Printf("  - %s (Value: $%d)\n", prop.Name, prop.Cost)
		totalValue += prop.Cost
	}
	fmt.Printf("Total Property Value: $%d\n", totalValue)
}

// GetRailroads returns all railroads owned by the player
func (p *Player) GetRailroads() []*Space {
	var railroads []*Space
	for _, prop := range p.Properties {
		if prop.Type == "Railroad" {
			railroads = append(railroads, prop)
		}
	}
	return railroads
}

// GetUtilities returns all utilities owned by the player
func (p *Player) GetUtilities() []*Space {
	var utilities []*Space
	for _, prop := range p.Properties {
		if prop.Type == "Utility" {
			utilities = append(utilities, prop)
		}
	}
	return utilities
}

// GetProperties returns all color properties owned by the player
func (p *Player) GetProperties() []*Space {
	var colorProps []*Space
	for _, prop := range p.Properties {
		if prop.Type == "Property" {
			colorProps = append(colorProps, prop)
		}
	}
	return colorProps
}

// PayMoney reduces player balance (for rent, taxes, etc.)
func (p *Player) PayMoney(amount int) bool {
	if p.Balance >= amount {
		p.Balance -= amount
		return true
	}
	return false
}

// ReceiveMoney increases player balance
func (p *Player) ReceiveMoney(amount int) {
	p.Balance += amount
}
