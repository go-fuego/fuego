package controller

import (
	"github.com/go-fuego/fuego"
)

type NewControllerRessources struct {
	// TODO add ressources
	NewControllerService NewControllerService
}

type NewController struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type NewControllerCreate struct {
	Name string `json:"name"`
}

type NewControllerUpdate struct {
	Name string `json:"name"`
}

func (rs NewControllerRessources) Routes(s *fuego.Server) {
	newControllerGroup := fuego.Group(s, "/newController")

	fuego.Get(newControllerGroup, "/", rs.getAllNewController)
	fuego.Post(newControllerGroup, "/", rs.postNewController)

	fuego.Get(newControllerGroup, "/{id}", rs.getNewController)
	fuego.Put(newControllerGroup, "/{id}", rs.putNewController)
	fuego.Delete(newControllerGroup, "/{id}", rs.deleteNewController)
}

func (rs NewControllerRessources) getAllNewController(c fuego.ContextNoBody) ([]NewController, error) {
	return rs.NewControllerService.GetAllNewController()
}

func (rs NewControllerRessources) postNewController(c *fuego.ContextWithBody[NewControllerCreate]) (NewController, error) {
	body, err := c.Body()
	if err != nil {
		return NewController{}, err
	}

	new, err := rs.NewControllerService.CreateNewController(body)
	if err != nil {
		return NewController{}, err
	}

	return new, nil
}

func (rs NewControllerRessources) getNewController(c fuego.ContextNoBody) (NewController, error) {
	id := c.PathParam("id")

	return rs.NewControllerService.GetNewController(id)
}

func (rs NewControllerRessources) putNewController(c *fuego.ContextWithBody[NewControllerUpdate]) (NewController, error) {
	id := c.PathParam("id")

	body, err := c.Body()
	if err != nil {
		return NewController{}, err
	}

	new, err := rs.NewControllerService.UpdateNewController(id, body)
	if err != nil {
		return NewController{}, err
	}

	return new, nil
}

func (rs NewControllerRessources) deleteNewController(c *fuego.ContextNoBody) (any, error) {
	return rs.NewControllerService.DeleteNewController(c.PathParam("id"))
}

type NewControllerService interface {
	GetNewController(id string) (NewController, error)
	CreateNewController(NewControllerCreate) (NewController, error)
	GetAllNewController() ([]NewController, error)
	UpdateNewController(id string, input NewControllerUpdate) (NewController, error)
	DeleteNewController(id string) (any, error)
}
