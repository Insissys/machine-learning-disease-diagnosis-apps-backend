package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/pkg/api"
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/pkg/model"
)

func Predict(c *gin.Context) {
	var request model.PredictRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, model.ApiResponse{Message: "Invalid request body"})
		return
	}

	translateRequest := &model.TranslateRequest{
		Q:      request.Symptoms,
		Source: "id",
		Target: "en",
		Format: "text",
	}

	translated, err := api.Translator(translateRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ApiResponse{
			Message: "Failed while translating symptoms",
			Error:   err,
		})
		return
	}

	request.Symptoms = translated.TranslatedText

	response, err := api.Predict(&request)
	if err != nil {
		if strings.Contains(err.Error(), "422") {
			c.JSON(http.StatusInternalServerError, model.ApiResponse{
				Message: "Can't predict: Try to be more verbose",
				Error:   err,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, model.ApiResponse{
			Message: "Something went wrong",
			Error:   err,
		})
		return
	}

	c.JSON(http.StatusOK, model.ApiResponse{
		Message: "Successfully Predict Disease",
		Data:    response,
	})
}
