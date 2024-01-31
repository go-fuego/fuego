package types

import "fmt"

type WhenToEat string

const (
	WhenToEatNone    WhenToEat = "-"
	WhenToEatStarter WhenToEat = "starter" // entrée
	WhenToEatDish    WhenToEat = "dish"    // plat
	WhenToEatDessert WhenToEat = "dessert" // dessert
)

// WhenToEatValues is a slice of all valid WhenToEats
var WhenToEatValues = []WhenToEat{
	WhenToEatNone,
	WhenToEatStarter,
	WhenToEatDish,
	WhenToEatDessert,
}

type InvalidWhenToEatError struct {
	WhenToEat WhenToEat
}

func (e InvalidWhenToEatError) Error() string {
	return fmt.Sprintf("invalid WhenToEat %s. Valid WhenToEats are: %v", e.WhenToEat, WhenToEatValues)
}

func (u WhenToEat) Valid() bool {
	for _, v := range WhenToEatValues {
		if v == u {
			return true
		}
	}
	return false
}
