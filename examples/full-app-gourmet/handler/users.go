package handler

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/examples/full-app-gourmet/store"
)

// MyCustomToken is a custom token that contains the standard claims and some custom claims.
type MyCustomToken struct {
	jwt.RegisteredClaims          // Required, this struct contains the standard claims
	Roles                []string `json:"roles"`
}

var _ jwt.Claims = &MyCustomToken{}

type LoginPayload struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

// Custom login controller
func (rs Resource) login(c fuego.ContextWithBody[LoginPayload]) (*TokenResponse, error) {
	body, err := c.Body()
	if err != nil {
		return nil, err
	}

	user, err := rs.UsersQueries.GetUserByUsername(context.Background(), body.Username)
	if err != nil {
		return nil, err
	}

	if user.EncryptedPassword != body.Password+"-encrypted" {
		return nil, fuego.UnauthorizedError{Title: "Unauthorized", Detail: "Invalid credentials"}
	}

	myToken := MyCustomToken{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    body.Username,
			Subject:   body.Username,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        "1234567890",
		},
		Roles: []string{"admin", "cook"},
	}

	s, err := rs.Security.GenerateToken(myToken)
	if err != nil {
		return nil, err
	}

	http.SetCookie(c.Response(), &http.Cookie{
		Name:     fuego.JWTCookieName,
		Value:    s,
		HttpOnly: true,
		MaxAge:   300,
		Domain:   os.Getenv("COOKIE_DOMAIN"),
		// SameSite: http.SameSiteStrictMode,
		// Secure:   true,
	})

	return &TokenResponse{
		Token: s,
	}, err
}

func (rs Resource) logout(c fuego.ContextNoBody) (any, error) {
	http.SetCookie(c.Response(), &http.Cookie{
		Name:   fuego.JWTCookieName,
		Value:  "",
		MaxAge: -1,
	})
	return nil, nil
}

func (rs Resource) me(c fuego.ContextNoBody) (*store.User, error) {
	caller, err := usernameFromContext(c.Context())
	if err != nil {
		return nil, err
	}

	user, err := rs.UsersQueries.GetUserByUsername(c.Context(), caller)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

type CreateUserPayload struct {
	Username string `json:"username" validate:"required"`
	FullName string `json:"full_name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

func (rs Resource) createUser(c fuego.ContextWithBody[CreateUserPayload]) (*store.User, error) {
	user, err := c.Body()
	if err != nil {
		return nil, err
	}

	createdUser, err := rs.UsersQueries.CreateUser(c.Context(), store.CreateUserParams{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		EncryptedPassword: user.Password + "-encrypted",
	})
	if err != nil {
		return nil, err
	}

	return &createdUser, nil
}

func (rs Resource) getUserByUsername(c fuego.ContextNoBody) (store.User, error) {
	username := c.PathParam("username")

	return rs.UsersQueries.GetUserByUsername(c.Context(), username)
}

func (rs Resource) getUsers(c fuego.ContextNoBody) ([]store.User, error) {
	return rs.UsersQueries.GetUsers(c.Context())
}

type UserRepository interface {
	CreateUser(ctx context.Context, arg store.CreateUserParams) (store.User, error)
	GetUserByUsername(ctx context.Context, username string) (store.User, error)
	GetUsers(ctx context.Context) ([]store.User, error)
}
