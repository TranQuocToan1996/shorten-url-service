package dto

type SubmitShortenURLRequest struct {
	LongURL string `json:"number" binding:"required,url"`
}

type GetDecodeURLRequest struct {
	ShortenURL string `json:"number" binding:"required,url"`
}
