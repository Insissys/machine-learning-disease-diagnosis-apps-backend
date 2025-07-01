package database

import (
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/internal/connection/db"
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/internal/database/migration"
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/pkg/model"
)

type DatabaseAuth struct{}

func NewDatabaseAuth() *DatabaseAuth {
	return &DatabaseAuth{}
}

func (a *DatabaseAuth) GetRefreshToken(id string) (*migration.Token, error) {
	var token *migration.Token

	err := db.Gorm.Debug().Where("id = ?", id).First(&token).Error
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (a *DatabaseAuth) StoreRefreshToken(token model.CustomClaims) error {
	data := &migration.Token{
		ID:      token.ID,
		User:    token.Email,
		Expired: token.ExpiresAt.Time,
		Revoked: token.Revoked,
	}
	err := db.Gorm.Debug().Create(&data).Error
	if err != nil {
		return err
	}

	return nil
}

func (a *DatabaseAuth) UpdateRefreshToken(token *migration.Token) error {
	err := db.Gorm.Debug().Where("id = ?", token.ID).Save(token).Error
	if err != nil {
		return err
	}

	return nil
}
