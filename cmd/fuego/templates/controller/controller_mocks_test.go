package controller

import "github.com/go-fuego/fuego"

type NewControllerServiceMock struct {
	getAllNewControllerLength int
}

var _ NewControllerService = NewControllerServiceMock{}

// CreateNewController implements NewControllerService.
func (NewControllerServiceMock) CreateNewController(body NewControllerCreate) (NewController, error) {
	return NewController{
		ID:   "randomID",
		Name: body.Name,
	}, nil
}

// DeleteNewController implements NewControllerService.
func (NewControllerServiceMock) DeleteNewController(id string) (any, error) {
	return nil, nil
}

// GetNewController implements NewControllerService.
func (NewControllerServiceMock) GetNewController(id string) (NewController, error) {
	if id == "404" {
		return NewController{}, fuego.HTTPError{
			Message:    "not found",
			StatusCode: 404,
		}
	}
	return NewController{}, nil
}

// UpdateNewController implements NewControllerService.
func (NewControllerServiceMock) UpdateNewController(id string, input NewControllerUpdate) (NewController, error) {
	return NewController{}, nil
}

func (m NewControllerServiceMock) GetAllNewController() ([]NewController, error) {
	allNewController := make([]NewController, m.getAllNewControllerLength)
	return allNewController, nil
}
