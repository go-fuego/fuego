package controller

import (
	"context"

	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/examples/full-app-gourmet/store"
)

type dosingRessource struct {
	Queries DosingRepository
}

func (rs dosingRessource) MountRoutes(s *fuego.Server) {
	fuego.Post(s, "/dosings/new", rs.newDosing)
}

func (rs dosingRessource) newDosing(c fuego.Ctx[store.CreateDosingParams]) (store.Dosing, error) {
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

type DosingRepository interface {
	CreateDosing(ctx context.Context, arg store.CreateDosingParams) (store.Dosing, error)
	GetDosings(ctx context.Context) ([]store.Dosing, error)
}

var _ DosingRepository = (*store.Queries)(nil)
