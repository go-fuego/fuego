package handlers

import (
	"strconv"

	"../crud-gorm/models"

	"crud-gorm/queries"

	"github.com/go-fuego/fuego"
)

type Handlers struct {
	UserQueries *queries.UserQueries
}

type UserToCreate struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (h *Handlers) GetUsers(c fuego.ContextNoBody) ([]models.User, error) {
	return h.UserQueries.GetUsers()
}

func (h *Handlers) GetUserByID(c fuego.ContextNoBody) (models.User, error) {
	id, err := strconv.Atoi(c.PathParam("id"))
	if err != nil {
		return models.User{}, fuego.BadRequestError{}
	}

	// Call UserQueries.GetUserByID which returns *models.User
	user, err := h.UserQueries.GetUserByID(uint(id))
	if err != nil {
		return models.User{}, err
	}

	// Dereference the pointer to return models.User
	return *user, nil
}

func (h *Handlers) UpdateUser(c fuego.ContextWithBody[models.User]) (models.User, error) {
	user := models.User{}
	error := h.UserQueries.UpdateUser(&user)
	return user, error
}

func (h *Handlers) CreateUser(c fuego.ContextWithBody[models.User]) (models.User, error) {
	user := models.User{}
	error := h.UserQueries.CreateUser(&user)
	return user, error
}

func (h *Handlers) DeleteUser(c fuego.ContextNoBody) (any, error) {
	id, error := strconv.Atoi(c.PathParam("id"))
	if error != nil {
		return nil, fuego.BadRequestError{}
	}
	error = h.UserQueries.DeleteUser(uint(id))
	return nil, error
}
