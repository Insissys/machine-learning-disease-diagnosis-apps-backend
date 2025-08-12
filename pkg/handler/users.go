package handler

import (
	"errors"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/internal/database/migration"
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/pkg/container"
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/pkg/model"
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/pkg/utils"
	"gorm.io/gorm"
)

// GET /api/users/me
func UsersMe(c *gin.Context) {
	email := c.MustGet("email").(string)
	db := container.NewContainer()

	user, err := db.Users.GetUserByEmail(&migration.User{Email: email})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, model.ApiResponse{Message: "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, model.ApiResponse{Message: "Something went wrong"})
		return
	}

	response := model.User{
		Base: model.Base{
			ID:        utils.EncryptUint64(uint64(user.ID)),
			CreatedAt: &user.CreatedAt,
			UpdatedAt: &user.UpdatedAt,
		},
		Name:  user.Name,
		Email: user.Email,
		Role: model.Roles{
			Name: user.Role.Name,
		},
		IsActive: user.IsActive,
		Expired:  user.Expired,
		Group: &model.Group{
			Base: model.Base{
				ID:        utils.EncryptUint64(uint64(user.Group.ID)),
				CreatedAt: &user.Group.CreatedAt,
				UpdatedAt: &user.Group.UpdatedAt,
			},
			Name:    user.Group.Name,
			Address: user.Group.Address,
		},
	}

	c.JSON(http.StatusOK, model.ApiResponse{
		Message:  "Succesfully Retrieve User",
		Error:    err,
		Data:     response,
		Metadata: nil,
	})
}

// GET /api/users
func GetUsers(c *gin.Context) {
	groupId := c.MustGet("groupId").(uint64)
	roles := c.Query("name")

	var rolesParam []string
	if roles == "" {
		rolesParam = append(rolesParam, []string{"admin", "doctor"}...)
	} else {
		re := regexp.MustCompile(`^[A-Za-z,]+$`)
		accepted := re.MatchString(roles)
		if !accepted {
			c.JSON(http.StatusBadRequest, model.ApiResponse{Message: "Invalid query"})
			return
		}

		splitted := strings.Split(roles, ",")

		rolesParam = append(rolesParam, splitted...)
	}

	db := container.NewContainer()

	users, err := db.Users.GetUsers(groupId, rolesParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ApiResponse{Message: "Something went wrong"})
		return
	}

	response := []model.User{}

	for _, v := range users {
		response = append(response, model.User{
			Base: model.Base{
				ID:        utils.EncryptUint64(uint64(v.ID)),
				CreatedAt: &v.CreatedAt,
				UpdatedAt: &v.UpdatedAt,
			},
			Name:  v.Name,
			Email: v.Email,
			Role: model.Roles{
				Name: v.Role.Name,
			},
			IsActive: v.IsActive,
			Expired:  v.Expired,
			Group: &model.Group{
				Base: model.Base{
					ID:        utils.EncryptUint64(uint64(v.Group.ID)),
					CreatedAt: &v.Group.CreatedAt,
					UpdatedAt: &v.Group.UpdatedAt,
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

// POST /api/users
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
		IsActive: user.IsActive,
		GroupID:  uint(c.MustGet("groupId").(uint64)),
	}
	if err := db.Users.StoreUser(data); err != nil {
		c.JSON(http.StatusInternalServerError, model.ApiResponse{Message: "Something went wrong"})
		return
	}

	c.JSON(http.StatusCreated, model.ApiResponse{
		Message: "User created successfully",
	})
}

// PATCH /api/users/:id
func PatchUser(c *gin.Context) {
	var user model.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, model.ApiResponse{Message: "Invalid request body"})
		return
	}

	encryptedBase64 := c.Param("id")
	id := utils.DecryptToUint64(encryptedBase64)
	db := container.NewContainer()

	if err := db.Users.PatchUser(id, &migration.User{
		Name:     user.Name,
		Email:    user.Email,
		IsActive: user.IsActive,
		Role: migration.Roles{
			Name: user.Role.Name,
		},
	}); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, model.ApiResponse{Message: "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, model.ApiResponse{Message: "Something went wrong"})
		return
	}

	c.JSON(http.StatusOK, model.ApiResponse{
		Message: "User updated successfully",
	})
}

// DELETE /api/users/:id
func DestroyUser(c *gin.Context) {
	encryptedBase64 := c.Param("id")
	id := utils.DecryptToUint64(encryptedBase64)
	db := container.NewContainer()

	if err := db.Users.DestroyUser(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, model.ApiResponse{Message: "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, model.ApiResponse{Message: "Something went wrong"})
		return
	}

	c.JSON(http.StatusOK, model.ApiResponse{
		Message: "User deleted successfully",
	})
}

// PATCH /api/users/activate/:id
func ActivateUser(c *gin.Context) {
	var data model.User

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, model.ApiResponse{Message: "Invalid request body"})
		return
	}

	encryptedBase64 := c.Param("id")
	id := utils.DecryptToUint64(encryptedBase64)
	db := container.NewContainer()

	if err := db.Users.ActivateUser(id, data.IsActive); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, model.ApiResponse{Message: "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, model.ApiResponse{Message: "Something went wrong"})
		return
	}

	c.JSON(http.StatusOK, model.ApiResponse{
		Message: "User activated successfully",
	})

}
