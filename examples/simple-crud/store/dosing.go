package store

import (
	"fmt"
	"simple-crud/store/types"
	"strings"

	"github.com/go-fuego/fuego"
)

var _ fuego.InTransformer = (*CreateDosingParams)(nil)

func (d *CreateDosingParams) InTransform() error {
	d.Unit = types.Unit(strings.ToLower(string(d.Unit)))

	if !d.Unit.Valid() {
		return types.InvalidUnitError{Unit: d.Unit}
	}

	if !(d.Quantity > 0 ||
		d.Quantity == 0 && d.Unit == types.UnitNone) {
		return fmt.Errorf("quantity must be greater than 0 for unit %s", d.Unit)
	}

	return nil
}
