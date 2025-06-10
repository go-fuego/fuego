# Hexagonal Architecture in Go with Fuego

Hexagonal Architecture, also known as Ports and Adapters Architecture, is a software design pattern that promotes separation of concerns by isolating the core business logic from external dependencies. This article demonstrates how to implement hexagonal architecture in Go using the Fuego web framework.

## What is Hexagonal Architecture?

Hexagonal Architecture was introduced by Alistair Cockburn to create loosely coupled application components that can be easily connected to their software environment through ports and adapters. The main benefits include:

- **Testability**: Business logic can be tested in isolation
- **Flexibility**: Easy to swap external dependencies
- **Maintainability**: Clear separation between business rules and infrastructure
- **Independence**: Core logic doesn't depend on frameworks or external systems

## Core Concepts

### The Hexagon (Core Domain)

The center contains your business logic, domain models, and use cases. This layer should have no dependencies on external frameworks or infrastructure.

### Ports

Ports are interfaces that define how the core domain communicates with the outside world. They come in two types:

- **Primary/Driving Ports**: Interfaces that allow external actors to interact with the application
- **Secondary/Driven Ports**: Interfaces that the application uses to interact with external systems

### Adapters

Adapters implement the ports and handle the actual communication with external systems:

- **Primary/Driving Adapters**: Handle incoming requests (HTTP handlers, CLI commands)
- **Secondary/Driven Adapters**: Handle outgoing requests (database repositories, external APIs)

## Project Structure

Let's organize our Go project following hexagonal architecture principles:

```
myapp/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── core/
│   │   ├── domain/
│   │   │   └── user.go
│   │   ├── ports/
│   │   │   ├── primary/
│   │   │   │   └── user_service.go
│   │   │   └── secondary/
│   │   │       └── user_repository.go
│   │   └── services/
│   │       └── user_service.go
│   └── adapters/
│       ├── primary/
│       │   └── http/
│       │       └── user_handler.go
│       └── secondary/
│           └── persistence/
│               └── user_repository.go
├── go.mod
└── go.sum
```

## Implementation Example

Let's build a simple user management system to demonstrate hexagonal architecture with Fuego.

### 1. Domain Layer

First, define your domain entities:

```go
// internal/core/domain/user.go
package domain

import (
    "errors"
    "context"
    "strings"
    "time"
)

type User struct {
    ID        string    `json:"id"`
    Email     string    `json:"email"`
    Name      string    `json:"name"`
    CreatedAt time.Time `json:"created_at"`
}

func NewUser(email, name string) (*User, error) {
    if email == "" {
        return nil, errors.New("email is required")
    }
    if name == "" {
        return nil, errors.New("name is required")
    }

    return &User{
        Email:     email,
        Name:      name,
        CreatedAt: time.Now(),
    }, nil
}

func (u *User) InTransform(_ context.Context) error {
    u.Email = strings.ToLower(u.Email)

    if u.Email == "" {
        return errors.New("email is required")
    }
    if u.Name == "" {
        return errors.New("name is required")
    }
    return nil
}
```

### 2. Ports (Interfaces)

Define the primary port (service interface):

```go
// internal/core/ports/primary/user_service.go
package primary

import (
    "context"
    "myapp/internal/core/domain"
)

type UserService interface {
    CreateUser(ctx context.Context, email, name string) (*domain.User, error)
    GetUser(ctx context.Context, id string) (*domain.User, error)
    GetAllUsers(ctx context.Context) ([]*domain.User, error)
    UpdateUser(ctx context.Context, id string, email, name string) (*domain.User, error)
    DeleteUser(ctx context.Context, id string) error
}
```

Define the secondary port (repository interface):

```go
// internal/core/ports/secondary/user_repository.go
package secondary

import (
    "context"
    "myapp/internal/core/domain"
)

type UserRepository interface {
    Save(ctx context.Context, user *domain.User) error
    FindByID(ctx context.Context, id string) (*domain.User, error)
    FindAll(ctx context.Context) ([]*domain.User, error)
    Update(ctx context.Context, user *domain.User) error
    Delete(ctx context.Context, id string) error
}
```

### 3. Core Service Implementation

Implement the business logic:

```go
// internal/core/services/user_service.go
package services

import (
    "context"
    "fmt"
    "myapp/internal/core/domain"
    "myapp/internal/core/ports/primary"
    "myapp/internal/core/ports/secondary"

    "github.com/google/uuid"
)

type userService struct {
    userRepo secondary.UserRepository
}

func NewUserService(userRepo secondary.UserRepository) primary.UserService {
    return &userService{
        userRepo: userRepo,
    }
}

func (s *userService) CreateUser(ctx context.Context, email, name string) (*domain.User, error) {
    user, err := domain.NewUser(email, name)
    if err != nil {
        return nil, fmt.Errorf("invalid user data: %w", err)
    }

    user.ID = uuid.New().String()

    if err := s.userRepo.Save(ctx, user); err != nil {
        return nil, fmt.Errorf("failed to save user: %w", err)
    }

    return user, nil
}

func (s *userService) GetUser(ctx context.Context, id string) (*domain.User, error) {
    user, err := s.userRepo.FindByID(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("failed to get user: %w", err)
    }
    return user, nil
}

func (s *userService) GetAllUsers(ctx context.Context) ([]*domain.User, error) {
    users, err := s.userRepo.FindAll(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed to get users: %w", err)
    }
    return users, nil
}

func (s *userService) UpdateUser(ctx context.Context, id string, email, name string) (*domain.User, error) {
    user, err := s.userRepo.FindByID(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("user not found: %w", err)
    }

    user.Email = email
    user.Name = name

    if err := user.Validate(); err != nil {
        return nil, fmt.Errorf("invalid user data: %w", err)
    }

    if err := s.userRepo.Update(ctx, user); err != nil {
        return nil, fmt.Errorf("failed to update user: %w", err)
    }

    return user, nil
}

func (s *userService) DeleteUser(ctx context.Context, id string) error {
    if err := s.userRepo.Delete(ctx, id); err != nil {
        return fmt.Errorf("failed to delete user: %w", err)
    }
    return nil
}
```

### 4. Secondary Adapter (Repository Implementation)

Implement the repository adapter:

```go
// internal/adapters/secondary/persistence/user_repository.go
package persistence

import (
    "context"
    "errors"
    "myapp/internal/core/domain"
    "myapp/internal/core/ports/secondary"
    "sync"
)

type inMemoryUserRepository struct {
    users map[string]*domain.User
    mutex sync.RWMutex
}

func NewInMemoryUserRepository() secondary.UserRepository {
    return &inMemoryUserRepository{
        users: make(map[string]*domain.User),
    }
}

func (r *inMemoryUserRepository) Save(ctx context.Context, user *domain.User) error {
    r.mutex.Lock()
    defer r.mutex.Unlock()

    r.users[user.ID] = user
    return nil
}

func (r *inMemoryUserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
    r.mutex.RLock()
    defer r.mutex.RUnlock()

    user, exists := r.users[id]
    if !exists {
        return nil, errors.New("user not found")
    }
    return user, nil
}

func (r *inMemoryUserRepository) FindAll(ctx context.Context) ([]*domain.User, error) {
    r.mutex.RLock()
    defer r.mutex.RUnlock()

    users := make([]*domain.User, 0, len(r.users))
    for _, user := range r.users {
        users = append(users, user)
    }
    return users, nil
}

func (r *inMemoryUserRepository) Update(ctx context.Context, user *domain.User) error {
    r.mutex.Lock()
    defer r.mutex.Unlock()

    if _, exists := r.users[user.ID]; !exists {
        return errors.New("user not found")
    }

    r.users[user.ID] = user
    return nil
}

func (r *inMemoryUserRepository) Delete(ctx context.Context, id string) error {
    r.mutex.Lock()
    defer r.mutex.Unlock()

    if _, exists := r.users[id]; !exists {
        return errors.New("user not found")
    }

    delete(r.users, id)
    return nil
}
```

### 5. Primary Adapter (HTTP Handler with Fuego)

Implement the HTTP adapter using Fuego:

```go
// internal/adapters/primary/http/user_handler.go
package http

import (
    "myapp/internal/core/domain"
    "myapp/internal/core/ports/primary"

    "github.com/go-fuego/fuego"
)

type UserHandler struct {
    userService primary.UserService
}

func NewUserHandler(userService primary.UserService) *UserHandler {
    return &UserHandler{
        userService: userService,
    }
}

type CreateUserRequest struct {
    Email string `json:"email" validate:"required,email"`
    Name  string `json:"name" validate:"required"`
}

type UpdateUserRequest struct {
    Email string `json:"email" validate:"required,email"`
    Name  string `json:"name" validate:"required"`
}

func (h *UserHandler) CreateUser(c fuego.ContextNoBody) (*domain.User, error) {
    var req CreateUserRequest
    if err := fuego.InBody(c, &req); err != nil {
        return nil, err
    }

    user, err := h.userService.CreateUser(c.Context(), req.Email, req.Name)
    if err != nil {
        return nil, err
    }

    return user, nil
}

func (h *UserHandler) GetUser(c fuego.ContextNoBody) (*domain.User, error) {
    id := c.PathParam("id")
    if id == "" {
        return nil, fuego.BadRequestError{Detail: "user ID is required"}
    }

    user, err := h.userService.GetUser(c.Context(), id)
    if err != nil {
        return nil, fuego.NotFoundError{Detail: err.Error()}
    }

    return user, nil
}

func (h *UserHandler) GetAllUsers(c fuego.ContextNoBody) ([]*domain.User, error) {
    users, err := h.userService.GetAllUsers(c.Context())
    if err != nil {
        return nil, err
    }

    return users, nil
}

func (h *UserHandler) UpdateUser(c fuego.ContextNoBody) (*domain.User, error) {
    id := c.PathParam("id")
    if id == "" {
        return nil, fuego.BadRequestError{Detail: "user ID is required"}
    }

    var req UpdateUserRequest
    if err := fuego.InBody(c, &req); err != nil {
        return nil, err
    }

    user, err := h.userService.UpdateUser(c.Context(), id, req.Email, req.Name)
    if err != nil {
        return nil, fuego.NotFoundError{Detail: err.Error()}
    }

    return user, nil
}

func (h *UserHandler) DeleteUser(c fuego.ContextNoBody) (any, error) {
    id := c.PathParam("id")
    if id == "" {
        return nil, fuego.BadRequestError{Detail: "user ID is required"}
    }

    if err := h.userService.DeleteUser(c.Context(), id); err != nil {
        return nil, fuego.NotFoundError{Detail: err.Error()}
    }

    return map[string]string{"message": "user deleted successfully"}, nil
}
```

### 6. Application Wiring

Finally, wire everything together in your main application:

```go
// cmd/server/main.go
package main

import (
    "log"
    "myapp/internal/adapters/primary/http"
    "myapp/internal/adapters/secondary/persistence"
    "myapp/internal/core/services"

    "github.com/go-fuego/fuego"
)

func main() {
    // Initialize repository (secondary adapter)
    userRepo := persistence.NewInMemoryUserRepository()

    // Initialize service (core business logic)
    userService := services.NewUserService(userRepo)

    // Initialize handler (primary adapter)
    userHandler := http.NewUserHandler(userService)

    // Setup Fuego server
    s := fuego.NewServer()

    // Register routes
    userGroup := fuego.Group(s, "/api/v1/users")
    fuego.Post(userGroup, "/", userHandler.CreateUser)
    fuego.Get(userGroup, "/", userHandler.GetAllUsers)
    fuego.Get(userGroup, "/{id}", userHandler.GetUser)
    fuego.Put(userGroup, "/{id}", userHandler.UpdateUser)
    fuego.Delete(userGroup, "/{id}", userHandler.DeleteUser)

    log.Println("Server starting on :8080")
    s.Run()
}
```

## Testing Hexagonal Architecture

One of the main benefits of hexagonal architecture is improved testability. You can easily test your business logic in isolation:

```go
// internal/core/services/user_service_test.go
package services_test

import (
    "context"
    "myapp/internal/core/domain"
    "myapp/internal/core/services"
    "testing"
)

// Mock repository for testing
type mockUserRepository struct {
    users map[string]*domain.User
}

func newMockUserRepository() *mockUserRepository {
    return &mockUserRepository{
        users: make(map[string]*domain.User),
    }
}

func (m *mockUserRepository) Save(ctx context.Context, user *domain.User) error {
    m.users[user.ID] = user
    return nil
}

func (m *mockUserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
    user, exists := m.users[id]
    if !exists {
        return nil, errors.New("user not found")
    }
    return user, nil
}

// ... implement other methods

func TestUserService_CreateUser(t *testing.T) {
    // Arrange
    mockRepo := newMockUserRepository()
    userService := services.NewUserService(mockRepo)

    // Act
    user, err := userService.CreateUser(context.Background(), "test@example.com", "Test User")

    // Assert
    if err != nil {
        t.Fatalf("expected no error, got %v", err)
    }

    if user.Email != "test@example.com" {
        t.Errorf("expected email 'test@example.com', got %s", user.Email)
    }

    if user.Name != "Test User" {
        t.Errorf("expected name 'Test User', got %s", user.Name)
    }

    if user.ID == "" {
        t.Error("expected user ID to be generated")
    }
}
```

## Benefits of This Architecture

### 1. Testability

Business logic can be tested without external dependencies. Mock implementations of repositories make unit testing straightforward.

### 2. Flexibility

You can easily swap the in-memory repository for a database implementation without changing the core business logic.

### 3. Framework Independence

The core domain doesn't depend on Fuego or any other framework. You could replace Fuego with another HTTP framework without affecting the business logic.

### 4. Clear Boundaries

The separation between layers makes it clear where different types of logic belong.

## Advanced Patterns

### Configuration and Environment Management

```go
// internal/config/config.go
package config

import (
    "os"
)

type Config struct {
    Port     string
    Database DatabaseConfig
}

type DatabaseConfig struct {
    Driver string
    DSN    string
}

func Load() *Config {
    return &Config{
        Port: getEnv("PORT", "8080"),
        Database: DatabaseConfig{
            Driver: getEnv("DB_DRIVER", "postgres"),
            DSN:    getEnv("DB_DSN", "postgres://localhost/myapp"),
        },
    }
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
```

### Error Handling

Create domain-specific errors:

```go
// internal/core/domain/errors.go
package domain

import "errors"

var (
    ErrUserNotFound      = errors.New("user not found")
    ErrUserAlreadyExists = errors.New("user already exists")
    ErrInvalidUserData   = errors.New("invalid user data")
)
```

### Middleware Integration

Fuego's middleware can be used at the adapter level without affecting the core:

```go
func setupMiddleware(s *fuego.Server) {
    fuego.Use(s, logger())
    fuego.Use(s, auth())
    // Add custom middleware for authentication, rate limiting, etc.
}
```

## Conclusion

Hexagonal Architecture provides a robust foundation for building maintainable Go applications. By using Fuego as the HTTP adapter, you get a powerful and type-safe web framework while keeping your business logic independent and testable.

Key takeaways:

- Keep your domain logic pure and framework-independent
- Use interfaces (ports) to define contracts between layers
- Implement adapters to handle external concerns
- Test your business logic in isolation using mocks
- Wire everything together at the application entry point

This architecture scales well as your application grows and makes it easier to maintain, test, and evolve your codebase over time.
