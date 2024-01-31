package store

import (
	"context"
	"strings"

	"github.com/go-fuego/fuego"
)

var _ fuego.InTransformer = (*CreateIngredientParams)(nil)

func (c *CreateIngredientParams) InTransform(context.Context) error {
	c.Name = strings.TrimSpace(c.Name)

	c.ID = slug(c.Name)

	return nil
}

func (i Ingredient) Months() string {
	months := []string{}
	if i.AvailableJan {
		months = append(months, "Jan")
	}
	if i.AvailableFeb {
		months = append(months, "Feb")
	}
	if i.AvailableMar {
		months = append(months, "Mar")
	}
	if i.AvailableApr {
		months = append(months, "Apr")
	}
	if i.AvailableMay {
		months = append(months, "May")
	}
	if i.AvailableJun {
		months = append(months, "Jun")
	}
	if i.AvailableJul {
		months = append(months, "Jul")
	}
	if i.AvailableAug {
		months = append(months, "Aug")
	}
	if i.AvailableSep {
		months = append(months, "Sep")
	}
	if i.AvailableOct {
		months = append(months, "Oct")
	}
	if i.AvailableNov {
		months = append(months, "Nov")
	}
	if i.AvailableDec {
		months = append(months, "Dec")
	}

	if len(months) == 0 {
		return "None"
	}

	return strings.Join(months, ", ")
}
