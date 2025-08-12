package model

type PredictRequest struct {
	Symptoms string `json:"symptoms" binding:"required"`
}

type DiseasePrediction struct {
	Disease    string  `json:"disease"`
	Confidence float64 `json:"confidence"`
}

type PredictResponse struct {
	Predictions     []DiseasePrediction `json:"predictions"`
	Recommendations []string            `json:"recommendations,omitempty"`
	Symptoms        []string            `json:"symptoms,omitempty"`
}
