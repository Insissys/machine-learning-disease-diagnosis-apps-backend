package database

import (
	"fmt"
	"time"

	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/internal/connection/db"
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/internal/database/migration"
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/pkg/utils"
	"gorm.io/gorm"
)

type DatabaseUsers struct{}

func NewDatabaseUsers() *DatabaseUsers {
	return &DatabaseUsers{}
}

func (d *DatabaseUsers) GetUserById(request *migration.User) (*migration.User, error) {
	db := db.Gorm
	response := &migration.User{}

	err := db.Preload("Group").Preload("Role").Where("id = ?", request.ID).First(&response).Error
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (d *DatabaseUsers) GetUserByEmail(request *migration.User) (*migration.User, error) {
	db := db.Gorm
	response := &migration.User{}

	err := db.Preload("Group").Preload("Role").Where("email = ?", request.Email).First(&response).Error
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (d *DatabaseUsers) RegisterUser(request *migration.User) error {
	return db.Gorm.Transaction(func(tx *gorm.DB) error {
		// 1. Create Group
		group := migration.Group{
			Name:    request.Group.Name,
			Address: request.Group.Address,
		}

		if err := tx.Create(&group).Error; err != nil {
			return fmt.Errorf("failed to create group: %w", err)
		}

		// 2. Get Role
		var role migration.Roles
		if err := tx.Where("name = ?", "superadmin").First(&role).Error; err != nil {
			return fmt.Errorf("failed to find role: %w", err)
		}

		// 3. Prepare User
		active := false
		user := &migration.User{
			Name:     request.Name,
			Email:    request.Email,
			Password: utils.HashPassword(request.Password),
			RoleID:   role.ID,
			Expired:  time.Now().Add(30 * 24 * time.Hour),
			GroupID:  group.ID,
			IsActive: active,
		}

		// 4. Save User
		if err := tx.Create(&user).Error; err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}

		return nil
	})
}

func (d *DatabaseUsers) GetUsers(groupId uint64, roles []string) ([]migration.User, error) {
	db := db.Gorm
	var response []migration.User

	err := db.
		Joins("JOIN roles ON roles.id = users.role_id").
		Where("users.group_id = ? AND roles.name IN (?)", groupId, roles).
		Preload("Group").
		Preload("Role").
		Find(&response).Error

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (d *DatabaseUsers) StoreUser(request *migration.User) error {
	db := db.Gorm

	var role migration.Roles
	err := db.Where("name = ?", request.Role.Name).First(&role).Error
	if err != nil {
		return err
	}

	data := &migration.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: utils.HashPassword(request.Password),
		IsActive: request.IsActive,
		RoleID:   role.ID,
		Expired:  time.Now().Add(30 * 24 * time.Hour),
		GroupID:  request.GroupID,
	}

	err = db.Create(&data).Error
	if err != nil {
		return err
	}

	return nil
}

func (d *DatabaseUsers) PatchUser(id uint64, data *migration.User) error {
	db := db.Gorm
	var user migration.User

	err := db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return err
	}

	var role migration.Roles
	err = db.Where("name = ?", data.Role.Name).First(&role).Error
	if err != nil {
		return err
	}

	u := migration.User{
		Name:     data.Name,
		Email:    data.Email,
		IsActive: data.IsActive,
		Role:     role,
	}

	err = db.Model(&user).Updates(&u).Error
	if err != nil {
		return err
	}

	return nil
}

func (d *DatabaseUsers) DestroyUser(id uint64) error {
	err := db.Gorm.Delete(&migration.User{}, id).Error
	if err != nil {
		return err
	}
	return nil
}

func (d *DatabaseUsers) ActivateUser(id uint64, isActive bool) error {
	db := db.Gorm
	var user migration.User

	err := db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return err
	}

	user.IsActive = isActive

	err = db.Save(&user).Error
	if err != nil {
		return err
	}

	return nil
}

func (d *DatabaseUsers) GetUserGroup(request *migration.User) (*migration.Group, error) {
	response := &migration.User{}

	err := db.Gorm.Preload("Group").Where("id = ?", request.ID).
		Or("email = ?", request.Email).
		Or("group_id = ?", request.GroupID).First(&response).Error
	if err != nil {
		return nil, err
	}
	return &response.Group, nil
}
