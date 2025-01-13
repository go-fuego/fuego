package fuego_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-fuego/fuego"
	"github.com/stretchr/testify/assert"
)

// UserProfile represents a user in our system
type UserProfile struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// UserService simulates a real service layer
type UserService interface {
	CreateUser(name, email string) (UserProfile, error)
	GetUserByID(id string) (UserProfile, error)
}

// mockUserService is a mock implementation of UserService
type mockUserService struct {
	users map[string]UserProfile
}

func newMockUserService() *mockUserService {
	return &mockUserService{
		users: map[string]UserProfile{
			"123": {ID: "123", Name: "John Doe", Email: "john@example.com"},
		},
	}
}

func (s *mockUserService) CreateUser(name, email string) (UserProfile, error) {
	if email == "taken@example.com" {
		return UserProfile{}, errors.New("email already taken")
	}
	user := UserProfile{
		ID:    "new_id",
		Name:  name,
		Email: email,
	}
	s.users[user.ID] = user
	return user, nil
}

func (s *mockUserService) GetUserByID(id string) (UserProfile, error) {
	user, exists := s.users[id]
	if !exists {
		return UserProfile{}, errors.New("user not found")
	}
	return user, nil
}

// CreateUserRequest represents the request body for user creation
type CreateUserRequest struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

// UserController handles user-related HTTP endpoints
type UserController struct {
	service UserService
}

func NewUserController(service UserService) *UserController {
	return &UserController{service: service}
}

// Create handles user creation
func (c *UserController) Create(ctx fuego.ContextWithBody[CreateUserRequest]) (UserProfile, error) {
	req, err := ctx.Body()
	if err != nil {
		return UserProfile{}, err
	}

	user, err := c.service.CreateUser(req.Name, req.Email)
	if err != nil {
		return UserProfile{}, err
	}

	ctx.SetStatus(http.StatusCreated)
	return user, nil
}

// GetByID handles fetching a user by ID
func (c *UserController) GetByID(ctx fuego.ContextWithBody[any]) (UserProfile, error) {
	id := ctx.PathParam("id")
	if id == "" {
		return UserProfile{}, errors.New("id is required")
	}

	user, err := c.service.GetUserByID(id)
	if err != nil {
		if err.Error() == "user not found" {
			ctx.SetStatus(http.StatusNotFound)
		}
		return UserProfile{}, err
	}

	return user, nil
}

func TestUserController(t *testing.T) {
	// Setup
	service := newMockUserService()
	controller := NewUserController(service)

	t.Run("create user", func(t *testing.T) {
		tests := []struct {
			name    string
			request CreateUserRequest
			want    UserProfile
			wantErr string
			status  int
		}{
			{
				name: "successful creation",
				request: CreateUserRequest{
					Name:  "Jane Doe",
					Email: "jane@example.com",
				},
				want: UserProfile{
					ID:    "new_id",
					Name:  "Jane Doe",
					Email: "jane@example.com",
				},
				status: http.StatusCreated,
			},
			{
				name: "email taken",
				request: CreateUserRequest{
					Name:  "Another User",
					Email: "taken@example.com",
				},
				wantErr: "email already taken",
				status:  http.StatusOK, // Default status when not set
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				// Setup mock context
				w := httptest.NewRecorder()
				ctx := fuego.NewMockContext[CreateUserRequest]()
				ctx.SetBody(tt.request)
				ctx.SetResponse(w)

				// Call controller
				got, err := controller.Create(ctx)

				// Assert results
				if tt.wantErr != "" {
					assert.Error(t, err)
					assert.Contains(t, err.Error(), tt.wantErr)
					return
				}

				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
				assert.Equal(t, tt.status, w.Code)
			})
		}
	})

	t.Run("get user by id", func(t *testing.T) {
		tests := []struct {
			name    string
			userID  string
			want    UserProfile
			wantErr string
			status  int
		}{
			{
				name:   "user found",
				userID: "123",
				want: UserProfile{
					ID:    "123",
					Name:  "John Doe",
					Email: "john@example.com",
				},
				status: http.StatusOK,
			},
			{
				name:    "user not found",
				userID:  "999",
				wantErr: "user not found",
				status:  http.StatusNotFound,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				// Setup mock context
				w := httptest.NewRecorder()
				ctx := fuego.NewMockContext[any]()
				ctx.SetPathParam("id", tt.userID)
				ctx.SetResponse(w)

				// Call controller
				got, err := controller.GetByID(ctx)

				// Assert results
				if tt.wantErr != "" {
					assert.Error(t, err)
					assert.Contains(t, err.Error(), tt.wantErr)
					assert.Equal(t, tt.status, w.Code)
					return
				}

				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
				assert.Equal(t, tt.status, w.Code)
			})
		}
	})
}
