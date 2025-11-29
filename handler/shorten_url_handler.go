package handler

import (
	"fmt"
	"net/http"

	"shorten/pkg/dto"
	"shorten/pkg/utils/url_utils"
	"shorten/service"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
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

func (h *ShortenURLHandler) SubmitEncode(c *gin.Context) {
	var req dto.SubmitShortenURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sendErrorResponse(c, http.StatusBadRequest, "fail", err.Error())
		return
	}
	err := h.urlService.SubmitURL(c.Request.Context(), req.LongURL)
	if err != nil {
		sendErrorResponse(c, http.StatusInternalServerError, "fail", err.Error())
		return
	}
	sendAPIResponse(c, http.StatusOK, "ok", "submitted", nil)
}

func (h *ShortenURLHandler) GetDecode(c *gin.Context) {
	var req dto.GetDecodeURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sendErrorResponse(c, http.StatusBadRequest, "fail", err.Error())
		return
	}
	urlObj, err := h.urlService.GetDecode(c.Request.Context(), req.ShortenURL)
	if err != nil {
		sendErrorResponse(c, http.StatusInternalServerError, "fail", err.Error())
		return
	}
	response := dto.GetDecodeURLResponse{}
	copier.Copy(urlObj, &response)
	sendAPIResponse(c, http.StatusOK, "ok", "success", response)
}

func (h *ShortenURLHandler) RedirectLongURL(c *gin.Context) {
	const maxCodeLength = 10
	code := c.Param("code")
	if !url_utils.IsBase62(code) || len(code) > maxCodeLength {
		sendErrorResponse(c, http.StatusBadRequest, "fail", fmt.Sprintf("code [%v] is not valid", code))
		return
	}
	urlObj, err := h.urlService.GetDecode(c.Request.Context(), code)
	if err != nil {
		sendErrorResponse(c, http.StatusInternalServerError, "fail", err.Error())
		return
	}
	c.Redirect(http.StatusMovedPermanently, urlObj.CleanURL)
}
