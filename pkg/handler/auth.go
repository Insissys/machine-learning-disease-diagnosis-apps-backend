package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/internal/config"
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/pkg/container"
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/pkg/model"
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/pkg/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Refresh(c *gin.Context) {
	cookie, err := c.Cookie("refresh_token")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Token are not found!"})
		return
	}

	cfg := config.GetConfig()

	// Parse JWT
	token, err := jwt.ParseWithClaims(cookie, &model.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return cfg.Config.JWTSECRET, nil
	})

	if err != nil || !token.Valid {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Invalid token"})
		return
	}

	claims, ok := token.Claims.(*model.CustomClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Unexpected while parsing token"})
		return
	}

	// Check token in store
	database := container.NewContainer()
	tokenDB, err := database.Auth.GetRefreshToken(claims.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Token not found"})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Unexpected while getting token from database"})
		return
	}

	if time.Now().After(tokenDB.Expired) {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Token expired"})
		return
	}

	if tokenDB.Revoked {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Token revoked"})
		return
	}

	// Rotate token: revoke old one
	tokenDB.Revoked = true
	err = database.Auth.UpdateRefreshToken(tokenDB)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Unexpected while updating token from database"})
		return
	}

	// Generate new tokens and store
	user, err := database.Users.GetUser(&model.LoginRequest{
		Email: tokenDB.User,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "User not found"})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Unexpected while getting user from database"})
		return
	}

	access, newRefresh, err := utils.GenerateTokens(model.User{
		Email:   user.Email,
		Role:    user.Role,
		GroupID: user.GroupID,
	}, database.Auth)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "token generation failed"})
		return
	}

	utils.SetRefreshCookie(c, newRefresh)
	c.JSON(http.StatusOK, gin.H{
		"token": access,
	})
}

func Login(c *gin.Context) {
	var input *model.LoginRequest = &model.LoginRequest{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	database := container.NewContainer()
	user, err := database.Users.GetUser(input)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	if !user.IsActive {
		c.JSON(http.StatusForbidden, gin.H{"message": "User inactive"})
		return
	}

	if time.Now().After(user.Expired) {
		c.JSON(http.StatusForbidden, gin.H{"message": "User expired"})
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid password"})
		return
	}

	// Generate JWT token
	access, refresh, err := utils.GenerateTokens(model.User{
		Email:   user.Email,
		Role:    user.Role,
		GroupID: user.GroupID,
	}, database.Auth)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "token generation failed"})
		return
	}

	utils.SetRefreshCookie(c, refresh)
	c.JSON(http.StatusOK, gin.H{
		"token":         access,
		"refresh_token": refresh,
	})
}

func Logout(c *gin.Context) {
	cookie, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Token are not found!"})
		return
	}

	cfg := config.GetConfig()

	// Parse JWT
	token, err := jwt.ParseWithClaims(cookie, &model.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return cfg.Config.JWTSECRET, nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
		return
	}

	claims, ok := token.Claims.(*model.CustomClaims)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Unexpected while parsing token"})
		return
	}

	// Check token in store
	database := container.NewContainer()
	tokenDB, err := database.Auth.GetRefreshToken(claims.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Token not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Unexpected while getting token from database"})
		return
	}

	if time.Now().After(tokenDB.Expired) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Token expired"})
		return
	}

	if tokenDB.Revoked {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Token revoked"})
		return
	}

	// Rotate token: revoke old one
	tokenDB.Revoked = true
	err = database.Auth.UpdateRefreshToken(tokenDB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Unexpected while updating token from database"})
		return
	}

	utils.RemoveCookie(c)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func Register(c *gin.Context) {
	var input *model.RegisterRequest = &model.RegisterRequest{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	database := container.NewContainer()
	err := database.Users.StoreUser(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created"})
}
