package controller

import (
	"simple-crud/store/dosings"

	"github.com/go-fuego/fuego"
)

type dosingRessource struct {
	Queries dosings.Queries
}

func (rs dosingRessource) MountRoutes(s *fuego.Server) {
	fuego.Post(s, "/dosings/new", rs.newDosing)
}

func (rs dosingRessource) newDosing(c fuego.Ctx[dosings.CreateDosingParams]) (dosings.Dosing, error) {
	body, err := c.Body()
	if err != nil {
		return dosings.Dosing{}, err
	}

	dosing, err := rs.Queries.CreateDosing(c.Context(), body)
	if err != nil {
		return dosings.Dosing{}, err
	}

	return dosing, nil
}
