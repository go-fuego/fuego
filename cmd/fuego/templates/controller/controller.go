package controller

import (
	"github.com/go-fuego/fuego"
)

type NewEntityResources struct {
	// TODO add resources
	NewEntityService NewEntityService
}

type NewEntity struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type NewEntityCreate struct {
	Name string `json:"name"`
}

type NewEntityUpdate struct {
	Name string `json:"name"`
}

func (rs NewEntityResources) Routes(s *fuego.Server) {
	newEntityGroup := fuego.Group(s, "/newEntity")

	fuego.Get(newEntityGroup, "/", rs.getAllNewEntity)
	fuego.Post(newEntityGroup, "/", rs.postNewEntity)

	fuego.Get(newEntityGroup, "/{id}", rs.getNewEntity)
	fuego.Put(newEntityGroup, "/{id}", rs.putNewEntity)
	fuego.Delete(newEntityGroup, "/{id}", rs.deleteNewEntity)
}

func (rs NewEntityResources) getAllNewEntity(c fuego.ContextNoBody) ([]NewEntity, error) {
	return rs.NewEntityService.GetAllNewEntity()
}

func (rs NewEntityResources) postNewEntity(c *fuego.ContextWithBody[NewEntityCreate]) (NewEntity, error) {
	body, err := c.Body()
	if err != nil {
		return NewEntity{}, err
	}

	new, err := rs.NewEntityService.CreateNewEntity(body)
	if err != nil {
		return NewEntity{}, err
	}

	return new, nil
}

func (rs NewEntityResources) getNewEntity(c fuego.ContextNoBody) (NewEntity, error) {
	id := c.PathParam("id")

	return rs.NewEntityService.GetNewEntity(id)
}

func (rs NewEntityResources) putNewEntity(c *fuego.ContextWithBody[NewEntityUpdate]) (NewEntity, error) {
	id := c.PathParam("id")

	body, err := c.Body()
	if err != nil {
		return NewEntity{}, err
	}

	new, err := rs.NewEntityService.UpdateNewEntity(id, body)
	if err != nil {
		return NewEntity{}, err
	}

	return new, nil
}

func (rs NewEntityResources) deleteNewEntity(c *fuego.ContextNoBody) (any, error) {
	return rs.NewEntityService.DeleteNewEntity(c.PathParam("id"))
}

type NewEntityService interface {
	GetNewEntity(id string) (NewEntity, error)
	CreateNewEntity(NewEntityCreate) (NewEntity, error)
	GetAllNewEntity() ([]NewEntity, error)
	UpdateNewEntity(id string, input NewEntityUpdate) (NewEntity, error)
	DeleteNewEntity(id string) (any, error)
}
