package database

import (
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/internal/connection/db"
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/internal/database/migration"
)

type DatabaseAuth struct{}

func NewDatabaseAuth() *DatabaseAuth {
	return &DatabaseAuth{}
}

func (a *DatabaseAuth) GetRefreshToken(id string) (*migration.Token, error) {
	var token *migration.Token

	err := db.Gorm.Where("id = ?", id).First(&token).Error
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (a *DatabaseAuth) StoreRefreshToken(token *migration.Token) error {
	data := &migration.Token{
		ID:      token.ID,
		User:    token.User,
		Expired: token.Expired,
		Revoked: token.Revoked,
	}
	err := db.Gorm.Create(&data).Error
	if err != nil {
		return err
	}

	return nil
}

func (a *DatabaseAuth) UpdateRefreshToken(token *migration.Token) error {
	err := db.Gorm.Where("id = ?", token.ID).Save(token).Error
	if err != nil {
		return err
	}

	return nil
}
