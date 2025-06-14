package api

import (
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/internal/config"
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/pkg/model"
)

func Predict(request *model.Request) (*model.Response, error) {
	cfg := config.GetConfig()
	client := resty.New()

	client.
		SetBaseURL(cfg.Config.Client.URL).
		SetTimeout(time.Duration(cfg.Config.Client.Timeout)*time.Second).
		SetRetryCount(3).
		SetHeader("Content-Type", "application/json")

	var response model.Response

	resp, err := client.R().
		SetBody(request).
		SetResult(&response).
		Post("/api/v1/predict")

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("fastapi error: %s", resp.Status())
	}

	return &response, nil
}
