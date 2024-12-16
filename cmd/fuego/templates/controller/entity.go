package controller

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

type NewEntityService interface {
	GetNewEntity(id string) (NewEntity, error)
	CreateNewEntity(NewEntityCreate) (NewEntity, error)
	GetAllNewEntity() ([]NewEntity, error)
	UpdateNewEntity(id string, input NewEntityUpdate) (NewEntity, error)
	DeleteNewEntity(id string) (any, error)
}
