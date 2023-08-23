package repository

import (
	"gorm.io/gorm"
	"log"
	"pkg/db/go/internal/models"
)

type UserRepository struct {
	dbClient *gorm.DB
}

func NewUserRepository(dbClient *gorm.DB) *UserRepository {
	return &UserRepository{
		dbClient: dbClient,
	}
}

func (repo *UserRepository) Insert(user *models.UserModel) error {
	err := repo.dbClient.Debug().
		Model(models.UserModel{}).
		Create(user).Error
	if err != nil {
		log.Printf("failed to insest user in database: %v\n", err)
		return err
	}

	return nil
}

func (repo *UserRepository) Delete(user *models.UserModel) error {
	err := repo.dbClient.Debug().Where("id = ?", user.ID).Delete(&models.UserModel{}).Error

	if err != nil {
		log.Print("failed to delete from db")
		return err
	}

	return nil
}

func (repo *UserRepository) CheckUserExists(userName string) (string, error) {
	var user models.UserModel
	result := repo.dbClient.Debug().Where("name= ?", userName).First(&user.Name)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return "", nil
		}
		return "", result.Error
	}

	return user.Password, nil
}
