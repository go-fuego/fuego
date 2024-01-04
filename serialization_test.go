package fuego

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

// crlf adds a crlf to the end of a string.
func crlf(s string) string {
	return s + "\n"
}

type response struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func TestRecursiveJSON(t *testing.T) {
	type rec struct {
		Rec *rec `json:"rec"`
	}

	t.Run("cannot serialize recursive json", func(t *testing.T) {
		w := httptest.NewRecorder()
		value := rec{}
		value.Rec = &value
		SendJSON(w, value)

		require.Equal(t, `{"error":"Cannot serialize JSON"}`, w.Body.String())
	})
}

func TestJSON(t *testing.T) {
	t.Run("can serialize json", func(t *testing.T) {
		w := httptest.NewRecorder()
		SendJSON(w, response{Message: "Hello World", Code: 200})
		body := w.Body.String()

		require.Equal(t, crlf(`{"message":"Hello World","code":200}`), body)
	})
}

func TestXML(t *testing.T) {
	t.Run("can serialize xml", func(t *testing.T) {
		w := httptest.NewRecorder()
		SendXML(w, response{Message: "Hello World", Code: 200})
		body := w.Body.String()

		require.Equal(t, `<response><Message>Hello World</Message><Code>200</Code></response>`, body)
	})

	t.Run("can serialize xml error", func(t *testing.T) {
		w := httptest.NewRecorder()
		err := HTTPError{Message: "Hello World"}
		SendXMLError(w, err)
		body := w.Body.String()

		require.Equal(t, `<HTTPError><Error>Hello World</Error></HTTPError>`, body)
	})
}

type tbt struct {
	Name string `json:"name"`
}

func (t *tbt) OutTransform(context.Context) error {
	t.Name = "transformed " + t.Name
	return nil
}

var _ OutTransformer = &tbt{}

func TestOutTranform(t *testing.T) {
	t.Run("can outTransform a value", func(t *testing.T) {
		value := tbt{Name: "John"}
		valueTransformed, err := transformOut(context.Background(), value)
		require.NoError(t, err)
		require.Equal(t, "transformed John", valueTransformed.Name)
	})

	t.Run("can outTransform a pointer to value", func(t *testing.T) {
		value := &tbt{Name: "Jack"}
		valueTransformed, err := transformOut(context.Background(), value)
		require.NoError(t, err)
		require.NotNil(t, valueTransformed)
		require.Equal(t, "transformed Jack", valueTransformed.Name)
	})

	t.Run("can outTransform a pointer to nil", func(t *testing.T) {
		valueTransformed, err := transformOut[*tbt](context.Background(), nil)
		require.NoError(t, err)
		require.Nil(t, valueTransformed)
	})

	t.Run("canNOT outTransform a value behind interface", func(t *testing.T) {
		value := tbt{Name: "Jack"}
		valueTransformed, err := transformOut[any](context.Background(), value)
		require.NoError(t, err)
		require.NotNil(t, valueTransformed)
		require.Equal(t, "Jack", valueTransformed.(tbt).Name)
	})

	t.Run("can outTransform a pointer to value behind interface", func(t *testing.T) {
		value := &tbt{Name: "Jack"}
		valueTransformed, err := transformOut[any](context.Background(), value)
		require.NoError(t, err)
		require.NotNil(t, valueTransformed)
		require.Equal(t, "transformed Jack", valueTransformed.(*tbt).Name)
	})
}

func BenchmarkOutTransform(b *testing.B) {
	b.Run("value", func(b *testing.B) {
		value := tbt{Name: "John"}
		for i := 0; i < b.N; i++ {
			value, err := transformOut(context.Background(), value)
			if err != nil {
				b.Fatal(err)
			}
			if value.Name != "transformed John" {
				b.Fatal("value not transformed")
			}
		}
	})

	b.Run("pointer to value", func(b *testing.B) {
		baseValue := tbt{Name: "Jack"}
		for i := 0; i < b.N; i++ {
			// Copy baseValue to value to avoid mutating the baseValue again and again.
			value := baseValue
			v, err := transformOut(context.Background(), &value)
			if err != nil {
				b.Fatal(err)
			}
			if v.Name != "transformed Jack" {
				b.Fatal("value not transformed on iteration", i, ": value", v)
			}
		}
	})

	b.Run("pointer to nil", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			value, err := transformOut[*tbt](context.Background(), nil)
			if err != nil {
				b.Fatal(err)
			}
			if value != nil {
				b.Fatal("value should be return nil")
			}
		}
	})
}

func TestJSONError(t *testing.T) {
	me := validatableStruct{
		Name:       "Napoleon Bonaparte",
		Age:        12,
		Email:      "not_an_email",
		ExternalID: "not_an_uuid",
	}

	err := validate(me)
	w := httptest.NewRecorder()
	err = ErrorHandler(err)
	SendJSONError(w, err)

	require.JSONEq(t, `
	{
		"error":"Name should be max=10, Age should be min=18, Required is required, Email should be a valid email, ExternalID should be a valid UUID",
		"info": {
		   "validation": [
			  {
				 "devField":"validatableStruct.Name",
				 "field":"Name",
				 "tag":"max",
				 "param":"10",
				 "value":"Napoleon Bonaparte"
			  },
			  {
				 "devField":"validatableStruct.Age",
				 "field":"Age",
				 "tag":"min",
				 "param":"18",
				 "value":12
			  },
			  {
				 "devField":"validatableStruct.Required",
				 "field":"Required",
				 "tag":"required",
				 "value":""
			  },
			  {
				 "devField":"validatableStruct.Email",
				 "field":"Email",
				 "tag":"email",
				 "value":"not_an_email"
			  },
			  {
				 "devField":"validatableStruct.ExternalID",
				 "field":"ExternalID",
				 "tag":"uuid",
				 "value":"not_an_uuid"
			  }
		   ]
		}
	 }`, w.Body.String())
}

func TestSend(t *testing.T) {
	w := httptest.NewRecorder()
	Send(w, "Hello World")

	require.Equal(t, "Hello World", w.Body.String())
}
