package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/pkg/api"
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/pkg/model"
)

func Predict(c *gin.Context) {
	var request *model.Request = &model.Request{}

	if err := c.ShouldBindJSON(request); err != nil {
		c.JSON(http.StatusInternalServerError, model.StatusInternalServerError(err.Error()))
		return
	}

	response, err := api.Predict(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.StatusInternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response)
}
