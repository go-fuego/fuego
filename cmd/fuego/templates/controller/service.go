package controller

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-fuego/fuego"
)

type NewControllerServiceImpl struct {
	newControllerRepository map[string]NewController
	mu                      sync.RWMutex
}

var _ NewControllerService = &NewControllerServiceImpl{}

func NewNewControllerService() NewControllerService {
	return &NewControllerServiceImpl{
		newControllerRepository: make(map[string]NewController),
	}
}

func (bs *NewControllerServiceImpl) GetNewController(id string) (NewController, error) {
	bs.mu.RLock()
	defer bs.mu.RUnlock()

	newController, exists := bs.newControllerRepository[id]
	if !exists {
		return NewController{}, fuego.NotFoundError{Title: "NewController not found with id " + id}
	}

	return newController, nil
}

func (bs *NewControllerServiceImpl) CreateNewController(input NewControllerCreate) (NewController, error) {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	id := fmt.Sprintf("%d", time.Now().UnixNano())
	newController := NewController{
		ID:   id,
		Name: input.Name,
	}

	bs.newControllerRepository[id] = newController
	return newController, nil
}

func (bs *NewControllerServiceImpl) GetAllNewController() ([]NewController, error) {
	bs.mu.RLock()
	defer bs.mu.RUnlock()

	newControllers := make([]NewController, 0, len(bs.newControllerRepository))
	for _, newController := range bs.newControllerRepository {
		newControllers = append(newControllers, newController)
	}

	return newControllers, nil
}

func (bs *NewControllerServiceImpl) UpdateNewController(id string, input NewControllerUpdate) (NewController, error) {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	newController, exists := bs.newControllerRepository[id]
	if !exists {
		return NewController{}, fuego.NotFoundError{Title: "NewController not found with id " + id}
	}

	if input.Name != "" {
		newController.Name = input.Name
	}

	bs.newControllerRepository[id] = newController
	return newController, nil
}

func (bs *NewControllerServiceImpl) DeleteNewController(id string) (any, error) {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	_, exists := bs.newControllerRepository[id]
	if !exists {
		return nil, fuego.NotFoundError{Title: "NewController not found with id " + id}
	}

	delete(bs.newControllerRepository, id)
	return "deleted", nil
}
