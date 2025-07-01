package utils

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/internal/config"
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/pkg/model"
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/pkg/repository"
)

func GenerateTokens(data model.User, repo repository.DatabaseAuthRepository) (string, string, error) {
	now := time.Now()
	cfg := config.GetConfig()

	// Access Token
	accessClaims := model.CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
		Email:   data.Email,
		Role:    data.Role,
		GroupID: data.GroupID,
	}
	accessJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err := accessJWT.SignedString(cfg.Config.JWTSECRET)
	if err != nil {
		return "", "", err
	}

	// Refresh Token
	jti := uuid.New().String()
	refreshClaims := model.CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        jti,
		},
		Email:   data.Email,
		Role:    data.Role,
		GroupID: data.GroupID,
		Revoked: false,
	}
	refreshJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err := refreshJWT.SignedString(cfg.Config.JWTSECRET)
	if err != nil {
		return "", "", err
	}

	// Store refresh token in DB
	err = repo.StoreRefreshToken(refreshClaims)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil

}

func SetRefreshCookie(c *gin.Context, token string) {
	c.SetCookie("refresh_token", token, int(24*time.Hour), "/", "", true, true)
}

func RemoveCookie(c *gin.Context) {
	c.SetCookie("refresh_token", "", -1, "/", "", true, true)
}
