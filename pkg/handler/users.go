package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/pkg/container"
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/pkg/model"
	"gorm.io/gorm"
)

func UsersMe(c *gin.Context) {
	email := c.MustGet("email").(string)
	db := container.NewContainer()

	user, err := db.Users.GetUser(&model.LoginRequest{Email: email})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": &model.User{
			Name:      user.Name,
			Email:     user.Email,
			Role:      user.Role,
			Expired:   user.Expired,
			GroupName: user.GroupName,
			Address:   user.Address,
		},
	})
}

func GetUsers(c *gin.Context) {
	groupId := c.MustGet("groupId").(uint)
	db := container.NewContainer()

	groupIdStr := strconv.Itoa(int(groupId))

	users, err := db.Users.GetUsers(groupIdStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}

func StoreUser(c *gin.Context) {
	var user struct {
		ID       uint   `json:"id"`
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Role     string `json:"role" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	db := container.NewContainer()

	data := model.User{
		Name:     user.Name,
		Email:    user.Email,
		Role:     user.Role,
		Password: user.Password,
		IsActive: false,                       // Default to active
		GroupID:  c.MustGet("groupId").(uint), // Get group ID
	}
	if err := db.Users.StoreUser(data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
	})
}

func PatchUser(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	id := c.Param("id")
	db := container.NewContainer()

	if err := db.Users.PatchUser(id, user); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
	})
}

func DestroyUser(c *gin.Context) {
	id := c.Param("id")
	db := container.NewContainer()

	if err := db.Users.DestroyUser(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
	})
}

func ActivateUser(c *gin.Context) {
	var data struct {
		IsActive bool `json:"isactive"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	id := c.Param("id")
	db := container.NewContainer()

	if err := db.Users.ActivateUser(id, data.IsActive); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User activated successfully",
	})

}
