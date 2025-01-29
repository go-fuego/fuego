package queries

import (
	"errors"

	"gorm.io/gorm"

	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/examples/crud-gorm/models"
)

type UserQueries struct {
	DB *gorm.DB
}

func (q *UserQueries) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	err := q.DB.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fuego.NotFoundError{
				Title:  "User not found",
				Detail: "No user with the provided ID was found.",
				Err:    err,
			}
		}
		return nil, err

	}
	return &user, nil
}

func (q *UserQueries) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := q.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (q *UserQueries) GetUsers() ([]models.User, error) {
	var users []models.User
	if err := q.DB.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (q *UserQueries) CreateUser(user *models.User) (*models.User, error) {
	err := q.DB.Create(user).Error
	if err != nil {
		return nil, fuego.InternalServerError{
			Detail: "Failed to create the user.",
			Err:    err,
		}
	}
	return user, nil
}

func (q *UserQueries) UpdateUser(user *models.User) (*models.User, error) {
	err := q.DB.Save(user).Error
	if err != nil {
		return nil, fuego.InternalServerError{
			Detail: "Failed to update the user.",
			Err:    err,
		}
	}
	return user, nil
}

func (q *UserQueries) DeleteUser(id uint) error {
	return q.DB.Delete(&models.User{}, id).Error
}
