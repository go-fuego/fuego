package fuego

import (
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

	var errStructValidation HTTPError
	require.ErrorAs(t, err, &errStructValidation)
	require.Equal(t, 400, errStructValidation.StatusCode())
	require.Equal(t, "Validation Error (400): Name should be max=10, Age should be min=18, Required is required, Email should be a valid email, ExternalID should be a valid UUID", errStructValidation.Error())
	require.Len(t, errStructValidation.Errors, 5)
}
