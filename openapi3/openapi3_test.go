package openapi3

import "testing"

func TestToSchema(t *testing.T) {

	t.Run("string", func(t *testing.T) {
		s := ToSchema("")
		if s.Type != "string" {
			t.Errorf("expected string, got %s", s.Type)
		}
	})

	t.Run("int", func(t *testing.T) {
		s := ToSchema(0)
		if s.Type != "object" {
			t.Errorf("expected object, got %s", s.Type)
		}
	})

	t.Run("struct", func(t *testing.T) {
		type S struct {
			A      string
			B      int
			Nested struct {
				C string
			}
		}
		s := ToSchema(S{})
		if s.Type != "object" {
			t.Errorf("expected object, got %s", s.Type)
		}
		if s.Properties["A"].Type != "string" {
			t.Errorf("expected string, got %s", s.Properties["A"].Type)
		}
		if s.Properties["B"].Type != "int" {
			t.Errorf("expected object, got %s", s.Properties["B"].Type)
		}
		if s.Properties["Nested"].Type != "object" {
			t.Errorf("expected object, got %s", s.Properties["Nested"].Type)
		}
	})

}
