package models

import (
	"context"
	"errors"

	"github.com/go-fuego/fuego"
)

type Pets struct {
	ID   string `json:"id" validate:"required" example:"pet-123456"`
	Name string `json:"name" validate:"required" example:"Napoleon"`
	Age  int    `json:"age" example:"18" description:"Age of the pet, in years"`
}

type PetsCreate struct {
	Name string `json:"name" validate:"required,min=1,max=100" example:"Napoleon"`
	Age  int    `json:"age" validate:"min=0,max=100" example:"18" description:"Age of the pet, in years"`
}

type PetsUpdate struct {
	Name string `json:"name,omitempty" validate:"min=1,max=100" example:"Napoleon" description:"Name of the pet"`
	Age  int    `json:"age,omitempty" validate:"max=100" example:"18"`
}

var _ fuego.InTransformer = &Pets{}

func (*Pets) InTransform(context.Context) error {
	return errors.New("pets must only be used as output")
}
