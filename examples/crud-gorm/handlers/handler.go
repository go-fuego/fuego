package handlers

import (
	"strconv"

	"gorm.io/gorm"

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
		if err == gorm.ErrRecordNotFound {
			return nil, fuego.NotFoundError{
				Title:  "User Not Found",
				Detail: "No user found with the provided ID.",
				Err:    err,
			}
		}
		return nil, fuego.InternalServerError{
			Detail: "An error occurred while checking the user.",
			Err:    err,
		}
	}

	input, err := c.Body()
	if err != nil {
		return nil, fuego.BadRequestError{
			Title:  "Invalid Input",
			Detail: "The provided input is not a valid user object.",
			Err:    err,
		}
	}

	existingUser.Name = input.Name
	existingUser.Email = input.Email

	updatedUser, err := h.UserQueries.UpdateUser(existingUser)
	if err != nil {
		return nil, fuego.InternalServerError{
			Detail: "Failed to update the user.",
			Err:    err,
		}
	}

	return updatedUser, nil
}

func (h *UserResources) CreateUser(c fuego.ContextWithBody[UserToCreate]) (*models.User, error) {
	input, err := c.Body()
	if err != nil {
		return nil, err
	}

	if input.Name == "" || input.Email == "" {
		return nil, fuego.BadRequestError{
			Title:  "Missing Required Fields",
			Detail: "The user name and email are required.",
			Err:    err,
		}
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
	createdUser, err := h.UserQueries.CreateUser(&userToCreate)
	if err != nil {
		return nil, fuego.ConflictError{
			Detail: "A user with the provided email already exists.",
			Err:    err,
		}
	}
	return createdUser, nil
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
		if err == gorm.ErrRecordNotFound {
			return err, fuego.NotFoundError{
				Title:  "User not found",
				Detail: "No user with the provided ID was found.",
				Err:    err,
			}
		}
		return err, fuego.InternalServerError{
			Detail: "An error occurred while retrieving the user.",
			Err:    err,
		}
	}

	err = h.UserQueries.DeleteUser(user.ID)
	if err != nil {
		return err, fuego.InternalServerError{
			Detail: "Failed to delete the user.",
			Err:    err,
		}
	}
	return nil, nil
}
