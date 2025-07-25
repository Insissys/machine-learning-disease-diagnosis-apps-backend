package model

type ApiResponse struct {
	Message  string `json:"message"`
	Error    error  `json:"error"`
	Data     any    `json:"data"`
	Metadata any    `json:"metadata"`
}
