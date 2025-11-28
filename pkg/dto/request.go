package dto

type APIResponse struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
