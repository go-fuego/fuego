package server

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"

	"github.com/go-fuego/fuego"
)

// TokenToContext is a middleware that checks if the user is authenticated from various authentication methods.
// Once found, the token is parsed, validated and the claims are set in the context.
// TLDR: after this middleware, the token is either non-existent or validated.
// You can use [TokenFromContext] to get the claims
func TokenToContext(security fuego.Security, searchFunc ...func(*http.Request) string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get the authorizationHeader from the header
			token := ""
			for _, f := range searchFunc {
				token = f(r)
				if token != "" {
					break
				}
			}

			if token == "" {
				// Unauthenticated, might be legit
				next.ServeHTTP(w, r)
				return
			}

			// Validate the token
			t, err := security.ValidateToken(token)
			if err != nil {
				// Remove cookie
				http.SetCookie(w, &http.Cookie{
					Name:   fuego.JWTCookieName,
					Domain: os.Getenv("COOKIE_DOMAIN"),
					Value:  "",
					MaxAge: -1,
				})
				fuego.SendError(w, r, fuego.UnauthorizedError{Title: "Unauthorized", Detail: "Invalid token", Err: err})
				return
			}

			// Get the claims
			claims := t.Claims.(jwt.MapClaims)

			// Set the subject and roles in the context
			ctx := r.Context()

			ctx = fuego.WithValue(ctx, claims)
			iss, _ := t.Claims.GetIssuer()
			ctx = context.WithValue(ctx, "issuer", iss)

			r = r.WithContext(ctx)

			// Call the next handler
			next.ServeHTTP(w, r)
		})
	}
}

func TokenFromCookie(r *http.Request) string {
	cookie, err := r.Cookie(fuego.JWTCookieName)
	if err != nil {
		return ""
	}

	if cookie == nil || cookie.Valid() != nil {
		slog.Info("Cookie is invalid", "cookie", cookie, "error", err)
		return ""
	}

	return cookie.Value
}
