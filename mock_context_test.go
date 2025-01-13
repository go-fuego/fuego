package fuego_test

import (
	"errors"
	"testing"

	"github.com/go-fuego/fuego"
	"github.com/stretchr/testify/assert"
)

// UserRequest represents the incoming request body
type UserRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// UserResponse represents the API response
type UserResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// CreateUserController is a typical controller that creates a user
func CreateUserController(c fuego.ContextWithBody[UserRequest]) (UserResponse, error) {
	// Get and validate the request body
	body, err := c.Body()
	if err != nil {
		return UserResponse{}, err
	}

	// Check if email is already taken (simulating DB check)
	if body.Email == "taken@example.com" {
		return UserResponse{}, errors.New("email already taken")
	}

	// In a real app, you would:
	// 1. Hash the password
	// 2. Save to database
	// 3. Generate ID
	// Here we'll simulate that:
	user := UserResponse{
		ID:    "user_123", // Simulated generated ID
		Name:  body.Name,
		Email: body.Email,
	}

	return user, nil
}

func TestCreateUserController(t *testing.T) {
	tests := []struct {
		name    string
		request UserRequest
		want    UserResponse
		wantErr string
	}{
		{
			name: "successful creation",
			request: UserRequest{
				Name:     "John Doe",
				Email:    "john@example.com",
				Password: "secure123",
			},
			want: UserResponse{
				ID:    "user_123",
				Name:  "John Doe",
				Email: "john@example.com",
			},
		},
		{
			name: "email already taken",
			request: UserRequest{
				Name:     "Jane Doe",
				Email:    "taken@example.com",
				Password: "secure123",
			},
			wantErr: "email already taken",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			ctx := fuego.NewMockContext[UserRequest]()
			ctx.SetBody(tt.request)

			// Execute
			got, err := CreateUserController(ctx)

			// Assert
			if tt.wantErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

// Example of testing a controller that uses path parameters
func GetUserController(c fuego.ContextNoBody) (UserResponse, error) {
	userID := c.PathParam("id")
	if userID == "" {
		return UserResponse{}, errors.New("user ID is required")
	}

	// Simulate fetching user from database
	if userID == "not_found" {
		return UserResponse{}, errors.New("user not found")
	}

	return UserResponse{
		ID:    userID,
		Name:  "John Doe",
		Email: "john@example.com",
	}, nil
}

func TestGetUserController(t *testing.T) {
	tests := []struct {
		name    string
		userID  string
		want    UserResponse
		wantErr string
	}{
		{
			name:   "user found",
			userID: "user_123",
			want: UserResponse{
				ID:    "user_123",
				Name:  "John Doe",
				Email: "john@example.com",
			},
		},
		{
			name:    "user not found",
			userID:  "not_found",
			wantErr: "user not found",
		},
		{
			name:    "missing user ID",
			userID:  "",
			wantErr: "user ID is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			ctx := fuego.NewMockContext[struct{}]()
			if tt.userID != "" {
				ctx.SetPathParam("id", tt.userID)
			}

			// Execute
			got, err := GetUserController(ctx)

			// Assert
			if tt.wantErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
} 