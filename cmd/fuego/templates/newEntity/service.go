package newEntity

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-fuego/fuego"
)

type NewEntityServiceImpl struct {
	newEntityRepository map[string]NewEntity
	mu                  sync.RWMutex
}

var _ NewEntityService = &NewEntityServiceImpl{}

func NewNewEntityService() NewEntityService {
	return &NewEntityServiceImpl{
		newEntityRepository: make(map[string]NewEntity),
	}
}

func (bs *NewEntityServiceImpl) GetNewEntity(id string) (NewEntity, error) {
	bs.mu.RLock()
	defer bs.mu.RUnlock()

	newEntity, exists := bs.newEntityRepository[id]
	if !exists {
		return NewEntity{}, fuego.NotFoundError{Title: "NewEntity not found with id " + id}
	}

	return newEntity, nil
}

func (bs *NewEntityServiceImpl) CreateNewEntity(input NewEntityCreate) (NewEntity, error) {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	id := fmt.Sprintf("%d", time.Now().UnixNano())
	newEntity := NewEntity{
		ID:   id,
		Name: input.Name,
	}

	bs.newEntityRepository[id] = newEntity
	return newEntity, nil
}

func (bs *NewEntityServiceImpl) GetAllNewEntity() ([]NewEntity, error) {
	bs.mu.RLock()
	defer bs.mu.RUnlock()

	allNewEntity := make([]NewEntity, 0, len(bs.newEntityRepository))
	for _, newEntity := range bs.newEntityRepository {
		allNewEntity = append(allNewEntity, newEntity)
	}

	return allNewEntity, nil
}

func (bs *NewEntityServiceImpl) UpdateNewEntity(id string, input NewEntityUpdate) (NewEntity, error) {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	newEntity, exists := bs.newEntityRepository[id]
	if !exists {
		return NewEntity{}, fuego.NotFoundError{Title: "NewEntity not found with id " + id}
	}

	if input.Name != "" {
		newEntity.Name = input.Name
	}

	bs.newEntityRepository[id] = newEntity
	return newEntity, nil
}

func (bs *NewEntityServiceImpl) DeleteNewEntity(id string) (any, error) {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	_, exists := bs.newEntityRepository[id]
	if !exists {
		return nil, fuego.NotFoundError{Title: "NewEntity not found with id " + id}
	}

	delete(bs.newEntityRepository, id)
	return "deleted", nil
}
