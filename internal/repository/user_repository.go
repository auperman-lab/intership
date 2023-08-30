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
		log.Printf("failed to insert user in database: %v\n", err)
		return err
	}

	return nil
}

func (repo *UserRepository) Delete(user *models.UserModel) error {
	err := repo.dbClient.Debug().Unscoped().Where("id = ?", user.ID).Delete(user).Error

	if err != nil {
		log.Print("failed to delete from postgres")
		return err
	}
	return nil

}

func (repo *UserRepository) CheckUserEmail(userName string) (string, uint, error) {
	var user models.UserModel
	result := repo.dbClient.Debug().Where("email= ?", userName).First(&user)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return "", 0, nil
		}
		return "", 0, result.Error
	}

	return user.Password, user.ID, nil
}

func (repo *UserRepository) CheckUserId(userID uint) (models.UserModel, error) {
	var user models.UserModel
	result := repo.dbClient.Debug().Where("id= ?", userID).First(&user)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return user, nil
		}
		return user, result.Error
	}

	return user, nil
}
