package fuego

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
	assert.Equal(t, 400, errStructValidation.Status)
	assert.Equal(t, "Validation Error", errStructValidation.Title)
	assert.Len(t, errStructValidation.Errors, 5)
	assert.Equal(t, "400 Validation Error (Name should be max=10, Age should be min=18, Required is required, Email should be a valid email, ExternalID should be a valid UUID)", errStructValidation.PublicError())
	assert.EqualError(t, errStructValidation, `400 Validation Error (Name should be max=10, Age should be min=18, Required is required, Email should be a valid email, ExternalID should be a valid UUID): Key: 'validatableStruct.Name' Error:Field validation for 'Name' failed on the 'max' tag
Key: 'validatableStruct.Age' Error:Field validation for 'Age' failed on the 'min' tag
Key: 'validatableStruct.Required' Error:Field validation for 'Required' failed on the 'required' tag
Key: 'validatableStruct.Email' Error:Field validation for 'Email' failed on the 'email' tag
Key: 'validatableStruct.ExternalID' Error:Field validation for 'ExternalID' failed on the 'uuid' tag`)
}
