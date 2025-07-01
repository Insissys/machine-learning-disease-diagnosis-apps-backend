package database

import (
	"time"

	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/internal/connection/db"
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/internal/database/migration"
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/pkg/model"
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/pkg/utils"
)

type DatabaseUsers struct{}

func NewDatabaseUsers() *DatabaseUsers {
	return &DatabaseUsers{}
}

func (d *DatabaseUsers) GetUser(request *model.LoginRequest) (*model.User, error) {
	db := db.Gorm
	var user migration.User

	err := db.Preload("Group").Preload("Role").Where("email = ?", request.Email).First(&user).Error
	if err != nil {
		return nil, err
	}

	response := model.User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role.Name,
		Password:  user.Password,
		IsActive:  *user.IsActive,
		GroupID:   user.Group.ID,
		Expired:   user.Expired,
		GroupName: user.Group.Name,
		Address:   user.Group.Address,
	}

	return &response, nil
}

func (d *DatabaseUsers) StoreUser(request *model.RegisterRequest) error {
	group := &migration.Group{
		Name:    request.GroupName,
		Address: request.Address,
	}

	err := db.Gorm.Debug().Create(&group).Error
	if err != nil {
		return err
	}

	var role migration.Roles
	err = db.Gorm.Debug().Where("name = ?", "superadmin").First(&role).Error
	if err != nil {
		return err
	}

	data := &migration.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: utils.HashPassword(request.Password),
		RoleID:   role.ID,
		Expired:  time.Now().Add(30 * 24 * time.Hour),
		GroupID:  group.ID,
	}

	active := false
	data.IsActive = &active

	err = db.Gorm.Debug().Create(&data).Error
	if err != nil {
		return err
	}

	return nil
}
