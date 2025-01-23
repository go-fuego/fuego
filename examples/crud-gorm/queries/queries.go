package queries

import (
	"github.com/sonkeydotcom/fuego/examples/crud-gorm/models"
	"gorm.io/gorm"
)

type UserQueries struct {
	DB *gorm.DB
}

func (q *UserQueries) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	err := q.DB.Where("id = ?", id).First(&user).Error
	return &user, err
}

func (q *UserQueries) GetUsers() ([]models.User, error) {
	var users []models.User
	if err := q.DB.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (q *UserQueries) CreateUser(user *models.User) error {
	return q.DB.Create(user).Error
}

func (q *UserQueries) UpdateUser(user *models.User) error {
	return q.DB.Save(user).Error
}

func (q *UserQueries) DeleteUser(id uint) error {
	return q.DB.Delete(&models.User{}, id).Error
}
