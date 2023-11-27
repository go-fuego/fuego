package views

import (
	"context"

	"simple-crud/store"
)

type DosingRepository interface {
	CreateDosing(ctx context.Context, arg store.CreateDosingParams) (store.Dosing, error)
	GetDosings(ctx context.Context) ([]store.Dosing, error)
}

var _ DosingRepository = (*store.Queries)(nil)
