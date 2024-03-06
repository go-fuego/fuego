package openapi3

import (
	"strconv"
	"strings"
)

// parses the "validate" tag used by go-playground/validator and sets the schema accordingly
func parseValidate(s *Schema, tag string) {
	for _, item := range strings.Split(tag, ",") {
		if strings.Contains(item, "min=") {
			min := strings.Split(item, "min=")[1]
			if s.Type == "integer" || s.Type == "number" {
				s.Minimum, _ = strconv.Atoi(min)
			} else if s.Type == "string" {
				s.MinLength, _ = strconv.Atoi(min)
			}
		}

		if strings.Contains(item, "max=") {
			max := strings.Split(item, "max=")[1]
			if s.Type == "integer" || s.Type == "number" {
				s.Maximum, _ = strconv.Atoi(max)
			} else if s.Type == "string" {
				s.MaxLength, _ = strconv.Atoi(max)
			}
		}
	}
}
