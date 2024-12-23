package fuego

import "fmt"

type ValidableCtx interface {
	GetOpenAPIParams() map[string]OpenAPIParam
	HasQueryParam(key string) bool
	HasHeader(key string) bool
	HasCookie(key string) bool
}

// ValidateParams checks if all required parameters are present in the request.
func ValidateParams(c ValidableCtx) error {
	for k, param := range c.GetOpenAPIParams() {
		if param.Default != nil {
			// skip: param has a default
			continue
		}

		if param.Required {
			switch param.Type {
			case QueryParamType:
				if !c.HasQueryParam(k) {
					err := fmt.Errorf("%s is a required query param", k)
					return BadRequestError{
						Title:  "Query Param Not Found",
						Err:    err,
						Detail: "cannot parse request parameter: " + err.Error(),
					}
				}
			case HeaderParamType:
				if !c.HasHeader(k) {
					err := fmt.Errorf("%s is a required header", k)
					return BadRequestError{
						Title:  "Header Not Found",
						Err:    err,
						Detail: "cannot parse request parameter: " + err.Error(),
					}
				}
			case CookieParamType:
				if !c.HasCookie(k) {
					err := fmt.Errorf("%s is a required cookie", k)
					return BadRequestError{
						Title:  "Cookie Not Found",
						Err:    err,
						Detail: "cannot parse request parameter: " + err.Error(),
					}
				}
			}
		}
	}

	return nil
}
