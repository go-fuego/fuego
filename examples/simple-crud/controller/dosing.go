package controller

import (
	"simple-crud/store"

	"github.com/go-fuego/fuego"
)

func (rs Ressource) newDosing(c fuego.Ctx[store.CreateDosingParams]) (store.Dosing, error) {
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
