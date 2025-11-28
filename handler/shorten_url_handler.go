package handler

import (
	"net/http"

	"shorten/pkg/dto"
	"shorten/service"

	"github.com/gin-gonic/gin"
)

type ShortenURLHandler struct {
	urlService service.URLService
}

func NewShortenURLHandler(
	urlService service.URLService,
) *ShortenURLHandler {
	return &ShortenURLHandler{
		urlService: urlService,
	}
}

func (h *ShortenURLHandler) SubmitShortURL(c *gin.Context) {
	var req dto.SubmitShortenURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sendErrorResponse(c, http.StatusBadRequest, "fail", err.Error())
		return
	}
	err := h.urlService.SubmitURL(c.Request.Context(), req.URL)
	if err != nil {
		sendErrorResponse(c, http.StatusInternalServerError, "fail", err.Error())
		return
	}
	sendAPIResponse(c, http.StatusOK, "ok", "submitted", nil)
}
