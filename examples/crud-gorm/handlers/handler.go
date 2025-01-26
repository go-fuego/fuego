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
		}
	}
	user, err := h.UserQueries.GetUserByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return models.User{}, fuego.NotFoundError{
				Title:  "User not found",
				Detail: "No user with the provided ID was found.",
			}
		}
		return models.User{}, fuego.InternalServerError{
			Detail: "An error occurred while retrieving the user.",
		}
	}
	return *user, nil
}

func (h *UserHandlers) UpdateUser(c fuego.ContextWithBody[models.User]) (models.User, error) {
	id, err := strconv.Atoi(c.PathParam("id"))

	if err != nil {
		return models.User{}, fuego.BadRequestError{
			Title:  "Inavalid ID",
			Detail: "The provided ID is not a valid integer.",
		}
	}

	input, err := c.Body()
	if err != nil {
		return models.User{}, fuego.BadRequestError{
			Title:  "Invalid Input",
			Detail: "The provided input is not a valid user object.",
		}
	}

	input.ID = uint(id)
	err = h.UserQueries.UpdateUser(&input)
	if err != nil {
		return models.User{}, fuego.InternalServerError{
			Detail: "Failed to update the user.",
		}
	}

	updatedUser, err := h.UserQueries.GetUserByID(uint(id))
	if err != nil {
		return models.User{}, fuego.InternalServerError{
			Detail: "Failed to retrieve the updated user.",
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
		}
	}

	if input.Name == "" || input.Email == "" {
		return models.User{}, fuego.BadRequestError{
			Title:  "Missing Required Fields",
			Detail: "The user name and email are required.",
		}
	}

	existingUser, err := h.UserQueries.GetUserByEmail(input.Email)
	if err == nil && existingUser != nil {
		return models.User{}, fuego.ConflictError{
			Detail: "A user with the provided email already exists.",
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
		}
	}

	user, err := h.UserQueries.GetUserByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return err, fuego.NotFoundError{
				Title:  "User not found",
				Detail: "No user with the provided ID was found.",
			}
		}
		return err, fuego.InternalServerError{
			Detail: "An error occurred while retrieving the user.",
		}
	}

	err = h.UserQueries.DeleteUser(user.ID)
	if err != nil {
		return err, fuego.InternalServerError{
			Detail: "Failed to delete the user.",
		}
	}
	return nil, nil
}
