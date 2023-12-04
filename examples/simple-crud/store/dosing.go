package store

import (
	"fmt"
	"strings"

	"github.com/go-fuego/fuego"
)

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

var _ fuego.InTransformer = (*CreateDosingParams)(nil)

func (d *CreateDosingParams) InTransform() error {
	d.Unit = strings.ToLower(string(d.Unit))

	if !Unit(d.Unit).Valid() {
		return InvalidUnitError{Unit: Unit(d.Unit)}
	}

	if !(d.Quantity > 0 ||
		d.Quantity == 0 && Unit(d.Unit) == UnitNone) {
		return fmt.Errorf("quantity must be greater than 0 for unit %s", d.Unit)
	}

	return nil
}
