package api

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/internal/config"
	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/pkg/model"
)

func Predict(request *model.PredictRequest) (*model.PredictResponse, error) {
	cfg := config.GetConfig()
	client := resty.New()

	client.
		SetBaseURL(cfg.Config.Model.URL).
		SetTimeout(time.Duration(cfg.Config.Model.Timeout)*time.Second).
		SetRetryCount(3).
		SetHeader("Content-Type", "application/json")

	response := &model.PredictResponse{}
	endpoint := "predict"

	resp, err := client.R().SetBody(request).SetResult(response).Post(endpoint)
	if err != nil {
		return nil, err
	}

	fmt.Println("HTTP Status:", resp.Status())
	fmt.Println("Response:", resp.String())

	if resp.StatusCode() == 422 {
		return nil, errors.New(resp.Status())
	}

	if resp.StatusCode() != 200 {
		return nil, errors.New(resp.String())
	}

	return response, nil
}

func Translator(request *model.TranslateRequest) (*model.TranslateResponse, error) {
	cfg := config.GetConfig()
	client := resty.New()

	client.
		SetBaseURL(cfg.Config.Translator.URL).
		SetTimeout(time.Duration(cfg.Config.Translator.Timeout)*time.Second).
		SetRetryCount(3).
		SetHeader("Content-Type", "application/json")

	response := &model.TranslateResponse{}
	endpoint := "translate"

	resp, err := client.R().SetBody(request).SetResult(response).Post(endpoint)
	if err != nil {
		return nil, err
	}

	fmt.Println("HTTP Status:", resp.Status())
	fmt.Println("Response:", resp.String())

	if resp.StatusCode() != 200 {
		return nil, errors.New(resp.String())
	}

	if response.TranslatedText == "" {
		return nil, errors.New("can not translate")
	}

	return response, nil
}
