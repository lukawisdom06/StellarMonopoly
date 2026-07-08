package game

import (
	"math/rand"
	"time"
)

// Dice represents a pair of dice
type Dice struct {
	Die1 int
	Die2 int
	Total int
}

// RollDice rolls two six-sided dice
func RollDice() Dice {
	rand.Seed(time.Now().UnixNano())
	die1 := rand.Intn(6) + 1
	die2 := rand.Intn(6) + 1
	
	return Dice{
		Die1:  die1,
		Die2:  die2,
		Total: die1 + die2,
	}
}

// IsDoubles checks if both dice show the same number
func (d Dice) IsDoubles() bool {
	return d.Die1 == d.Die2
}
