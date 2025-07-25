package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/internal/config"
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/internal/database/migration"
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/pkg/container"
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/pkg/model"
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/pkg/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// POST /api/refresh-token
func Refresh(c *gin.Context) {
	cookie, err := c.Cookie("refresh_token")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.ApiResponse{Message: "Token are not found!"})
		return
	}

	cfg := config.GetConfig()

	// Parse JWT
	token, err := jwt.ParseWithClaims(cookie, &model.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return cfg.Config.JWTSECRET, nil
	})

	if err != nil || !token.Valid {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.ApiResponse{Message: "Invalid token"})
		return
	}

	claims, ok := token.Claims.(*model.CustomClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.ApiResponse{Message: "Unexpected while parsing token"})
		return
	}

	// Check token in store
	database := container.NewContainer()
	tokenDB, err := database.Auth.GetRefreshToken(claims.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusInternalServerError, model.ApiResponse{Message: "Token not found"})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.ApiResponse{Message: "Unexpected while getting token from database"})
		return
	}

	if time.Now().After(tokenDB.Expired) {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.ApiResponse{Message: "Token expired"})
		return
	}

	if tokenDB.Revoked {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.ApiResponse{Message: "Token revoked"})
		return
	}

	// Rotate token: revoke old one
	tokenDB.Revoked = true
	err = database.Auth.UpdateRefreshToken(tokenDB)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.ApiResponse{Message: "Unexpected while updating token from database"})
		return
	}

	// Generate new tokens and store
	user, err := database.Users.GetUserByEmail(&migration.User{
		Email: tokenDB.User,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusInternalServerError, model.ApiResponse{Message: "User not found"})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.ApiResponse{Message: "Unexpected while getting user from database"})
		return
	}

	access, newRefresh, err := utils.GenerateTokens(model.User{
		Email: user.Email,
		Role:  model.Roles{Name: user.Role.Name},
		Group: model.Group{
			Base: model.Base{
				ID: user.GroupID,
			},
		},
	}, database.Auth)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.ApiResponse{Message: "token generation failed"})
		return
	}

	utils.SetRefreshCookie(c, newRefresh)
	c.JSON(http.StatusOK, model.ApiResponse{
		Message: "Success refreshing token",
		Data: gin.H{
			"token": access,
		},
	})
}

// POST /api/login
func Login(c *gin.Context) {
	var input model.LoginRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, model.ApiResponse{Message: "Invalid input"})
		return
	}

	database := container.NewContainer()
	user, err := database.Users.GetUserByEmail(&migration.User{
		Email: input.Email,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, model.ApiResponse{Message: "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, model.ApiResponse{Message: "Something went wrong", Error: err})
		return
	}

	if user.IsActive == nil || !*user.IsActive {
		c.JSON(http.StatusForbidden, model.ApiResponse{Message: "User inactive"})
		return
	}

	if time.Now().After(user.Expired) {
		c.JSON(http.StatusForbidden, model.ApiResponse{Message: "User expired"})
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, model.ApiResponse{Message: "Invalid password"})
		return
	}

	access, refresh, err := utils.GenerateTokens(model.User{
		Email: user.Email,
		Role: model.Roles{
			Name: user.Role.Name,
		},
		Group: model.Group{
			Base: model.Base{
				ID: user.GroupID,
			},
		},
	}, database.Auth)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ApiResponse{Message: "Token generation failed"})
		return
	}

	utils.SetRefreshCookie(c, refresh)

	c.JSON(http.StatusOK, model.ApiResponse{
		Message: "Successfully login",
		Data: gin.H{
			"token":         access,
			"refresh_token": refresh,
		},
	})
}

// POST /api/logout
func Logout(c *gin.Context) {
	cookie, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, model.ApiResponse{Message: "Token are not found!"})
		return
	}

	cfg := config.GetConfig()

	// Parse JWT
	token, err := jwt.ParseWithClaims(cookie, &model.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return cfg.Config.JWTSECRET, nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, model.ApiResponse{Message: "Invalid token"})
		return
	}

	claims, ok := token.Claims.(*model.CustomClaims)
	if !ok {
		c.JSON(http.StatusInternalServerError, model.ApiResponse{Message: "Unexpected while parsing token"})
		return
	}

	// Check token in store
	database := container.NewContainer()
	tokenDB, err := database.Auth.GetRefreshToken(claims.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusUnauthorized, model.ApiResponse{Message: "Token not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, model.ApiResponse{Message: "Unexpected while getting token from database"})
		return
	}

	if time.Now().After(tokenDB.Expired) {
		c.JSON(http.StatusUnauthorized, model.ApiResponse{Message: "Token expired"})
		return
	}

	if tokenDB.Revoked {
		c.JSON(http.StatusUnauthorized, model.ApiResponse{Message: "Token revoked"})
		return
	}

	// Rotate token: revoke old one
	tokenDB.Revoked = true
	err = database.Auth.UpdateRefreshToken(tokenDB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ApiResponse{Message: "Unexpected while updating token from database"})
		return
	}

	utils.RemoveCookie(c)
	c.JSON(http.StatusOK, model.ApiResponse{Message: "Logged out successfully"})
}

// POST /api/register
func Register(c *gin.Context) {
	var input model.RegisterRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, model.ApiResponse{Message: "Invalid input", Error: err})
		return
	}

	database := container.NewContainer()

	err := database.Users.RegisterUser(&migration.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
		Group: migration.Group{
			Name:    input.OfficeName,
			Address: input.Address,
		},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ApiResponse{Message: "Failed to register user", Error: err})
		return
	}

	c.JSON(http.StatusOK, model.ApiResponse{Message: "User registered successfully"})
}
