package model

type Request struct {
	Complaint string `json:"complaint" binding:"required"` // binding for Gin validation
}

type DiseasePrediction struct {
	Disease    string  `json:"disease"`    // predicted disease name
	Confidence float64 `json:"confidence"` // confidence score (0-1)
}

type Response struct {
	Predictions []DiseasePrediction `json:"predictions"`
}

type RegisterRequest struct {
	OfficeName string `json:"officename" binding:"required"`
	Address    string `json:"address" binding:"required"`
	Name       string `json:"name" binding:"required"`
	Email      string `json:"email" binding:"required"`
	Password   string `json:"password" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
