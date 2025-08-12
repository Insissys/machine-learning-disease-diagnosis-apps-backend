package model

type TranslateRequest struct {
	Q            string `json:"q"`
	Source       string `json:"source"`
	Target       string `json:"target"`
	Format       string `json:"format"`
	Alternatives string `json:"alternatives,omitempty"`
	ApiKey       string `json:"api_key,omitempty"`
}

type TranslateResponse struct {
	Alternatives   []string `json:"alternatives,omitempty"`
	TranslatedText string   `json:"translatedText"`
}
