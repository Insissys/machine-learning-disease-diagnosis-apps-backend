package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/internal/database/migration"
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/pkg/container"
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/pkg/model"
	"gorm.io/gorm"
)

func UsersMe(c *gin.Context) {
	email := c.MustGet("email").(string)
	db := container.NewContainer()

	user, err := db.Users.GetUserByEmail(&migration.User{Email: email})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, model.ApiResponse{Message: "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, model.ApiResponse{Message: "Something went wrong", Error: err})
		return
	}

	c.JSON(http.StatusOK, model.ApiResponse{
		Message: "Succesfully Retrieve User",
		Error:   err,
		Data: model.User{
			Base: model.Base{
				ID:        user.ID,
				CreatedAt: user.CreatedAt,
				UpdatedAt: user.UpdatedAt,
			},
			Name:     user.Name,
			Email:    user.Email,
			Password: user.Password,
			Role: model.Roles{
				Base: model.Base{
					ID:        user.Role.ID,
					CreatedAt: user.Role.CreatedAt,
					UpdatedAt: user.Role.UpdatedAt,
				},
				Name: user.Role.Name,
			},
			IsActive: *user.IsActive,
			Expired:  user.Expired,
			Group: model.Group{
				Base: model.Base{
					ID:        user.Group.ID,
					CreatedAt: user.Group.CreatedAt,
					UpdatedAt: user.Group.UpdatedAt,
				},
				Name:    user.Group.Name,
				Address: user.Group.Address,
			},
		},
		Metadata: nil,
	})
}

func GetUsers(c *gin.Context) {
	groupId := c.MustGet("groupId").(uint)
	db := container.NewContainer()

	groupIdStr := strconv.Itoa(int(groupId))

	users, err := db.Users.GetUsers(groupIdStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ApiResponse{Message: "Something went wrong", Error: err})
		return
	}

	response := []model.User{}

	for _, v := range users {
		response = append(response, model.User{
			Base: model.Base{
				ID:        v.ID,
				CreatedAt: v.CreatedAt,
				UpdatedAt: v.UpdatedAt,
			},
			Name:     v.Name,
			Email:    v.Email,
			Password: v.Password,
			Role: model.Roles{
				Base: model.Base{
					ID:        v.Role.ID,
					CreatedAt: v.Role.CreatedAt,
					UpdatedAt: v.Role.UpdatedAt,
				},
				Name: v.Role.Name,
			},
			IsActive: *v.IsActive,
			Expired:  v.Expired,
			Group: model.Group{
				Base: model.Base{
					ID:        v.Group.ID,
					CreatedAt: v.Group.CreatedAt,
					UpdatedAt: v.Group.UpdatedAt,
				},
				Name:    v.Group.Name,
				Address: v.Group.Address,
			},
		})
	}

	c.JSON(http.StatusOK, model.ApiResponse{
		Message:  "Succesfully Retrieve Users",
		Error:    err,
		Data:     response,
		Metadata: nil,
	})
}

func StoreUser(c *gin.Context) {
	var user model.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, model.ApiResponse{Message: "Invalid request body"})
		return
	}

	db := container.NewContainer()

	data := &migration.User{
		Name:  user.Name,
		Email: user.Email,
		Role: migration.Roles{
			Name: user.Role.Name,
		},
		Password: user.Password,
		IsActive: &user.IsActive,
		Group: migration.Group{
			Model: gorm.Model{
				ID: c.MustGet("groupId").(uint),
			},
		},
	}
	if err := db.Users.StoreUser(data); err != nil {
		c.JSON(http.StatusInternalServerError, model.ApiResponse{Message: "Something went wrong", Error: err})
		return
	}

	c.JSON(http.StatusCreated, model.ApiResponse{
		Message: "User created successfully",
	})
}

func PatchUser(c *gin.Context) {
	var user model.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, model.ApiResponse{Message: "Invalid request body"})
		return
	}

	id := c.Param("id")
	db := container.NewContainer()

	if err := db.Users.PatchUser(id, &migration.User{
		Name:     user.Name,
		Email:    user.Email,
		IsActive: &user.IsActive,
		Role: migration.Roles{
			Name: user.Role.Name,
		},
	}); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, model.ApiResponse{Message: "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, model.ApiResponse{Message: "Something went wrong", Error: err})
		return
	}

	c.JSON(http.StatusOK, model.ApiResponse{
		Message: "User updated successfully",
	})
}

func DestroyUser(c *gin.Context) {
	id := c.Param("id")
	db := container.NewContainer()

	if err := db.Users.DestroyUser(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, model.ApiResponse{Message: "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, model.ApiResponse{Message: "Something went wrong", Error: err})
		return
	}

	c.JSON(http.StatusOK, model.ApiResponse{
		Message: "User deleted successfully",
	})
}

func ActivateUser(c *gin.Context) {
	var data model.User

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, model.ApiResponse{Message: "Invalid request body"})
		return
	}

	id := c.Param("id")
	db := container.NewContainer()

	if err := db.Users.ActivateUser(id, data.IsActive); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, model.ApiResponse{Message: "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, model.ApiResponse{Message: "Something went wrong", Error: err})
		return
	}

	c.JSON(http.StatusOK, model.ApiResponse{
		Message: "User activated successfully",
	})

}
