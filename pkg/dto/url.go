package dto

type SubmitShortenURLRequest struct {
	URL string `json:"number" binding:"required"`
}
