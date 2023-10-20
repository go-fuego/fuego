package controller

import (
	"simple-crud/store"

	"github.com/go-op/op"
)

func (rs Ressource) newDosing(c op.Ctx[store.CreateDosingParams]) (store.Dosing, error) {
	body, err := c.Body()
	if err != nil {
		return store.Dosing{}, err
	}

	dosing, err := rs.Queries.CreateDosing(c.Context(), body)
	if err != nil {
		return store.Dosing{}, err
	}

	return dosing, nil
}
