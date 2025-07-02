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

func (d *DatabaseUsers) RegisterUser(request *model.RegisterRequest) error {
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

func (d *DatabaseUsers) GetUsers(groupId string) ([]model.User, error) {
	db := db.Gorm
	var users []migration.User

	err := db.Preload("Group").Preload("Role").Where("group_id = ?", groupId).Find(&users).Error
	if err != nil {
		return nil, err
	}

	var response []model.User
	for _, user := range users {
		if user.Role.Name == "superadmin" {
			continue // Skip superadmin users
		}
		response = append(response, model.User{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			Role:      user.Role.Name,
			GroupID:   user.Group.ID,
			Expired:   user.Expired,
			GroupName: user.Group.Name,
			Address:   user.Group.Address,
			IsActive:  *user.IsActive,
		})
	}

	return response, nil
}

func (d *DatabaseUsers) StoreUser(request model.User) error {
	db := db.Gorm

	var role migration.Roles
	err := db.Where("name = ?", request.Role).First(&role).Error
	if err != nil {
		return err
	}
	data := &migration.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: utils.HashPassword(request.Password),
		IsActive: &request.IsActive,
		RoleID:   role.ID,
		Expired:  time.Now().Add(30 * 24 * time.Hour),
		GroupID:  request.GroupID,
	}

	err = db.Debug().Create(&data).Error
	if err != nil {
		return err
	}

	return nil
}

func (d *DatabaseUsers) PatchUser(id string, data model.User) error {
	db := db.Gorm
	var user migration.User

	err := db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return err
	}

	u := migration.User{
		Name:     data.Name,
		Email:    data.Email,
		IsActive: &data.IsActive,
		Role:     migration.Roles{Name: data.Role},
		GroupID:  data.GroupID,
	}

	err = db.Model(&user).Updates(&u).Error
	if err != nil {
		return err
	}

	return nil
}

func (d *DatabaseUsers) DestroyUser(id string) error {
	err := db.Gorm.Debug().Delete(&migration.User{}, id).Error
	if err != nil {
		return err
	}
	return nil
}

func (d *DatabaseUsers) ActivateUser(id string, isActive bool) error {
	db := db.Gorm
	var user migration.User

	err := db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return err
	}

	user.IsActive = &isActive

	err = db.Save(&user).Error
	if err != nil {
		return err
	}

	return nil
}
