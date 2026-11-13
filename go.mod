module github.com/go-fuego/fuego

go 1.26.10

require (
	github.com/getkin/kin-openapi v0.153.0
	github.com/go-playground/validator/v10 v10.30.6
	github.com/golang-jwt/jwt/v5 v5.3.1
	github.com/google/uuid v1.6.0
	github.com/gorilla/schema v1.4.1
	github.com/stretchr/testify v1.12.1
	github.com/thejerf/slogassert v0.3.4
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/gabriel-vasile/mimetype v1.4.16 // indirect
	github.com/go-openapi/jsonpointer v1.0.1 // indirect
	github.com/go-playground/locales v0.14.2 // indirect
	github.com/go-playground/universal-translator v0.18.2 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/leodido/go-urn v1.5.0 // indirect
	github.com/rogpeppe/go-internal v1.16.0 // indirect
	github.com/santhosh-tekuri/jsonschema/v6 v6.0.3 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	go.yaml.in/yaml/v4 v4.0.0 // indirect
	golang.org/x/crypto v0.58.0 // indirect
	golang.org/x/sys v0.50.0 // indirect
	golang.org/x/text v0.44.0 // indirect
)

retract (
	v1.0.1 // Contains retractions only.
	v1.0.0 // Published accidentally.
)
