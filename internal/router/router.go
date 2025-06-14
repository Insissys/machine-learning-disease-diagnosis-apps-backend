package router

import (
	"github.com/gin-gonic/gin"
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/pkg/handler"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/predict", handler.Predict)
	// Add more routes here

	return r
}
