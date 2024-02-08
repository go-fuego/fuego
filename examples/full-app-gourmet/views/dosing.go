package views

import (
	"context"

	"github.com/go-fuego/fuego/examples/full-app-gourmet/store"
)

type DosingRepository interface {
	CreateDosing(ctx context.Context, arg store.CreateDosingParams) (store.Dosing, error)
	GetDosings(ctx context.Context) ([]store.Dosing, error)
}

var _ DosingRepository = (*store.Queries)(nil)
