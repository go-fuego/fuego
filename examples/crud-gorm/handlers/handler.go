package handlers

import (
	"strconv"

	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/examples/crud-gorm/models"
)

// UserResources is used to inject rependencies into the handlers about the user queries.
type UserResources struct {
	UserQueries UserQueryInterface
}

type UserToCreate struct {
	Name  string `json:"name" validate:"required,min=3"`
	Email string `json:"email" validate:"required,email"`
}

type UserQueryInterface interface {
	GetUsers() ([]models.User, error)
	GetUserByID(id uint) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	CreateUser(user *models.User) (*models.User, error)
	UpdateUser(user *models.User) (*models.User, error)
	DeleteUser(id uint) error
}

func (h *UserResources) GetUsers(c fuego.ContextNoBody) ([]models.User, error) {
	return h.UserQueries.GetUsers()
}

func (h *UserResources) GetUserByID(c fuego.ContextNoBody) (*models.User, error) {
	id, err := strconv.Atoi(c.PathParam("id"))
	if err != nil {
		return nil, fuego.BadRequestError{
			Title:  "Invalid ID",
			Detail: "The provided ID is not a valid integer.",
			Err:    err,
		}
	}

	return h.UserQueries.GetUserByID(uint(id))
}

func (h *UserResources) UpdateUser(c fuego.ContextWithBody[models.User]) (*models.User, error) {
	id, err := strconv.Atoi(c.PathParam("id"))
	if err != nil {
		return nil, fuego.BadRequestError{
			Title:  "Invalid ID",
			Detail: "The provided ID is not a valid integer.",
			Err:    err,
		}
	}

	existingUser, err := h.UserQueries.GetUserByID(uint(id))
	if err != nil {
		return nil, err
	}

	input, err := c.Body()
	if err != nil {
		return nil, err
	}

	existingUser.Name = input.Name
	existingUser.Email = input.Email

	return h.UserQueries.UpdateUser(existingUser)
}

func (h *UserResources) CreateUser(c fuego.ContextWithBody[UserToCreate]) (*models.User, error) {
	input, err := c.Body()
	if err != nil {
		return nil, err
	}

	existingUser, err := h.UserQueries.GetUserByEmail(input.Email)
	if err == nil && existingUser != nil {
		return nil, fuego.ConflictError{
			Detail: "A user with the provided email already exists.",
			Err:    err,
		}
	}

	userToCreate := models.User{
		Name:  input.Name,
		Email: input.Email,
	}

	return h.UserQueries.CreateUser(&userToCreate)
}

func (h *UserResources) DeleteUser(c fuego.ContextNoBody) (any, error) {
	id, err := strconv.Atoi(c.PathParam("id"))
	if err != nil {
		return err, fuego.BadRequestError{
			Title:  "Invalid ID",
			Detail: "The provided ID is not a valid integer.",
			Err:    err,
		}
	}

	user, err := h.UserQueries.GetUserByID(uint(id))
	if err != nil {
		return nil, err
	}

	err = h.UserQueries.DeleteUser(user.ID)

	return nil, err

}
