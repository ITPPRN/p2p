package repository

import (
	"gorm.io/gorm"

	"p2p-back-end/modules/entities/models"
)

type authRepositoryDB struct {
	db *gorm.DB
}

func NewAuthRepositoryDB(db *gorm.DB) models.UserRepository {
	return &authRepositoryDB{db: db}
}

func (r authRepositoryDB) IsUserExistByID(id string) (bool, error) {

	var count int64
	if err := r.db.Table("user_entities").Where("id = ?", id).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
