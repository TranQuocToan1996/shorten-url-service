package dto

type CalculateRequest struct {
	URL string `json:"number" binding:"required"`
}
