package handlers

import (
	"strconv"

	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/examples/crud-gorm/models"
	"gorm.io/gorm"
)

type UserHandlers struct {
	UserQueries UserQueryInterface
}

type UserToCreate struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserQueryInterface interface {
	GetUsers() ([]models.User, error)
	GetUserByID(id uint) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	CreateUser(user *models.User) error
	UpdateUser(user *models.User) error
	DeleteUser(id uint) error
}

func (h *UserHandlers) GetUsers(c fuego.ContextNoBody) ([]models.User, error) {
	return h.UserQueries.GetUsers()
}

func (h *UserHandlers) GetUserByID(c fuego.ContextNoBody) (models.User, error) {
	id, err := strconv.Atoi(c.PathParam("id"))
	if err != nil {
		return models.User{}, fuego.BadRequestError{
			Title:  "Invalid ID",
			Detail: "The provided ID is not a valid integer.",
			Err:    err,
		}
	}
	user, err := h.UserQueries.GetUserByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return models.User{}, fuego.NotFoundError{
				Title:  "User not found",
				Detail: "No user with the provided ID was found.",
				Err:    err,
			}
		}
		return models.User{}, fuego.InternalServerError{
			Detail: "An error occurred while retrieving the user.",
			Err:    err,
		}
	}
	return *user, nil
}

func (h *UserHandlers) UpdateUser(c fuego.ContextWithBody[models.User]) (models.User, error) {
	id, err := strconv.Atoi(c.PathParam("id"))
	if err != nil {
		return models.User{}, fuego.BadRequestError{
			Title:  "Invalid ID",
			Detail: "The provided ID is not a valid integer.",
			Err:    err,
		}
	}

	existingUser, err := h.UserQueries.GetUserByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return models.User{}, fuego.NotFoundError{
				Title:  "User Not Found",
				Detail: "No user found with the provided ID.",
				Err:    err,
			}
		}
		return models.User{}, fuego.InternalServerError{
			Detail: "An error occurred while checking the user.",
			Err:    err,
		}
	}

	input, err := c.Body()
	if err != nil {
		return models.User{}, fuego.BadRequestError{
			Title:  "Invalid Input",
			Detail: "The provided input is not a valid user object.",
			Err:    err,
		}
	}

	existingUser.Name = input.Name
	existingUser.Email = input.Email

	err = h.UserQueries.UpdateUser(existingUser)
	if err != nil {
		return models.User{}, fuego.InternalServerError{
			Detail: "Failed to update the user.",
			Err:    err,
		}
	}

	updatedUser, err := h.UserQueries.GetUserByID(uint(id))
	if err != nil {
		return models.User{}, fuego.InternalServerError{
			Detail: "Failed to retrieve the updated user.",
			Err:    err,
		}
	}

	return *updatedUser, nil
}

func (h *UserHandlers) CreateUser(c fuego.ContextWithBody[models.User]) (models.User, error) {
	input, err := c.Body()
	if err != nil {
		return models.User{}, fuego.BadRequestError{
			Title:  "Invalid Input",
			Detail: "The provided input is not a valid user object.",
			Err:    err,
		}
	}

	if input.Name == "" || input.Email == "" {
		return models.User{}, fuego.BadRequestError{
			Title:  "Missing Required Fields",
			Detail: "The user name and email are required.",
			Err:    err,
		}
	}

	existingUser, err := h.UserQueries.GetUserByEmail(input.Email)
	if err == nil && existingUser != nil {
		return models.User{}, fuego.ConflictError{
			Detail: "A user with the provided email already exists.",
			Err:    err,
		}
	}

	user := models.User{
		Name:  input.Name,
		Email: input.Email,
	}
	err = h.UserQueries.CreateUser(&user)
	if err != nil {
		return models.User{}, fuego.ConflictError{
			Detail: "A user with the provided email already exists.",
			Err:    err,
		}
	}
	return user, nil

}

func (h *UserHandlers) DeleteUser(c fuego.ContextNoBody) (any, error) {
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
