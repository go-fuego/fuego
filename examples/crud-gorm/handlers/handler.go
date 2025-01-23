package handlers

import (
	"strconv"

	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/examples/crud-gorm/models"
	"github.com/go-fuego/fuego/examples/crud-gorm/queries"
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
	user, err := h.UserQueries.GetUserByID(uint(id))
	if err != nil {
		return models.User{}, err
	}
	return *user, nil
}

func (h *Handlers) UpdateUser(c fuego.ContextWithBody[models.User]) (models.User, error) {
	id, err := strconv.Atoi(c.PathParam("id"))
	if err != nil {
		return models.User{}, fuego.BadRequestError{}
	}
	input, err := c.Body()
	if err != nil {
		return models.User{}, fuego.BadRequestError{}
	}
	input.ID = uint(id)
	err = h.UserQueries.UpdateUser(&input)
	if err != nil {
		return models.User{}, err
	}
	return input, nil
}

func (h *Handlers) CreateUser(c fuego.ContextWithBody[models.User]) (models.User, error) {
	input, err := c.Body()
	if err != nil {
		return models.User{}, fuego.BadRequestError{}
	}

	user := models.User{
		Name:  input.Name,
		Email: input.Email,
	}
	err = h.UserQueries.CreateUser(&user)
	return user, err

}

func (h *Handlers) DeleteUser(c fuego.ContextNoBody) (any, error) {
	id, error := strconv.Atoi(c.PathParam("id"))
	if error != nil {
		return nil, fuego.BadRequestError{}
	}
	error = h.UserQueries.DeleteUser(uint(id))
	return nil, error
}
