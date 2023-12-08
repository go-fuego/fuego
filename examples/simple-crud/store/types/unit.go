package types

import "fmt"

type Unit string

const (
	UnitNone       Unit = "-" // for ingredients that are not dosed, like salt, pepper, ...
	UnitPiece      Unit = "piece"
	UnitGramm      Unit = "g"
	UnitMilliliter Unit = "ml"
)

// UnitValues is a slice of all valid units
var UnitValues = []Unit{
	UnitNone,
	UnitPiece,
	UnitGramm,
	UnitMilliliter,
}

type InvalidUnitError struct {
	Unit Unit
}

func (e InvalidUnitError) Error() string {
	return fmt.Sprintf("invalid unit %s. Valid units are: %v", e.Unit, UnitValues)
}

func (u Unit) Valid() bool {
	for _, v := range UnitValues {
		if v == u {
			return true
		}
	}
	return false
}
