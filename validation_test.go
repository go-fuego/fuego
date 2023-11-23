package fuego

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type validatableStruct struct {
	Name       string `validate:"required,min=3,max=10"`
	Age        int    `validate:"min=18"`
	Required   string `validate:"required"`
	Email      string `validate:"email"`
	ExternalID string `validate:"uuid"`
}

func TestValidate(t *testing.T) {
	me := validatableStruct{
		Name:  "Napoleon Bonaparte",
		Age:   12,
		Email: "napoleon.bonaparte",
	}

	err := validate(me)
	t.Log(err)
	require.Error(t, err)

	var errStructValidation structValidationError
	if errors.As(err, &errStructValidation) {
		require.Equal(t, errStructValidation.Status(), 400)
		require.Equal(t, errStructValidation.Error(), "Name should be max=10, Age should be min=18, Required is required, Email should be a valid email, ExternalID should be a valid UUID")
		require.Len(t, errStructValidation.Errors, 5)
	} else {
		t.Error("error is not a structValidationError but should be")
	}
}
