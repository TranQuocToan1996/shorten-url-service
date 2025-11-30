package handler

import (
	"fmt"
	"net/http"
	"strings"

	"shorten/pkg/config"
	"shorten/pkg/dto"
	"shorten/pkg/utils/url_utils"
	"shorten/service"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

type ShortenURLHandler struct {
	urlService service.URLService
	config     config.Config
}

func NewShortenURLHandler(
	urlService service.URLService,
	config config.Config,
) *ShortenURLHandler {
	return &ShortenURLHandler{
		urlService: urlService,
		config:     config,
	}
}

// SubmitEncode godoc
// @Summary      Submit URL for encoding
// @Description  Submit a long URL to be shortened. The encoding will be processed asynchronously.
// @Tags         url
// @Accept       json
// @Produce      json
// @Param        request  body      dto.SubmitShortenURLRequest  true  "URL encoding request"
// @Success      200      {object}  dto.APIResponse
// @Failure      400      {object}  dto.APIResponse
// @Failure      500      {object}  dto.APIResponse
// @Router       /api/v1/encode [post]
func (h *ShortenURLHandler) SubmitEncode(c *gin.Context) {
	var req dto.SubmitShortenURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sendErrorResponse(c, http.StatusBadRequest, "fail", err.Error())
		return
	}
	err := h.urlService.SubmitURL(c.Request.Context(), req.LongURL, req.CallbackURL)
	if err != nil {
		sendErrorResponse(c, http.StatusInternalServerError, "fail", err.Error())
		return
	}
	sendAPIResponse(c, http.StatusOK, "ok", "submitted", nil)
}

// GetDecode godoc
// @Summary      Decode shortened URL
// @Description  Get the original long URL from a shortened URL
// @Tags         url
// @Accept       json
// @Produce      json
// @Param        shorten_url  query     string  true  "Shortened URL"
// @Success      200          {object}  dto.APIResponse{data=dto.GetDecodeURLResponse}
// @Failure      400          {object}  dto.APIResponse
// @Failure      500          {object}  dto.APIResponse
// @Router       /api/v1/decode [get]
func (h *ShortenURLHandler) GetDecode(c *gin.Context) {
	var req dto.GetDecodeURLRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		sendErrorResponse(c, http.StatusBadRequest, "fail", err.Error())
		return
	}
	code := strings.TrimPrefix(req.ShortenURL, h.config.REDIRECT_HOST)
	urlObj, err := h.urlService.GetDecode(c.Request.Context(), code)
	if err != nil {
		sendErrorResponse(c, http.StatusInternalServerError, "fail", err.Error())
		return
	}
	response := dto.GetDecodeURLResponse{}
	copier.Copy(&response, urlObj)
	sendAPIResponse(c, http.StatusOK, "ok", "success", response)
}

// GetURLEncodeByLongURL godoc
// @Summary      Get encoded URL by long URL
// @Description  Retrieve the shortened URL information by providing the original long URL
// @Tags         url
// @Accept       json
// @Produce      json
// @Param        long_url  query     string  true  "Original long URL"
// @Success      200       {object}  dto.APIResponse{data=dto.GetDecodeURLResponse}
// @Failure      400       {object}  dto.APIResponse
// @Failure      500       {object}  dto.APIResponse
// @Router       /api/v1/urls/long [get]
func (h *ShortenURLHandler) GetURLEncodeByLongURL(c *gin.Context) {
	var req dto.GetEncodeURLRequestByLongURL
	if err := c.ShouldBindQuery(&req); err != nil {
		sendErrorResponse(c, http.StatusBadRequest, "fail", err.Error())
		return
	}
	urlObj, err := h.urlService.GetByLongURL(c.Request.Context(), req.LongURL)
	if err != nil {
		sendErrorResponse(c, http.StatusInternalServerError, "fail", err.Error())
		return
	}
	response := dto.GetDecodeURLResponse{}
	copier.Copy(&response, urlObj)
	sendAPIResponse(c, http.StatusOK, "ok", "success", response)
}

// RedirectLongURL godoc
// @Summary      Redirect to long URL
// @Description  Redirect to the original long URL using the shortened code
// @Tags         url
// @Produce      json
// @Param        code  path      string  true  "Shortened URL code"
// @Success      301   {string}  string  "Redirect to long URL"
// @Failure      400   {object}  dto.APIResponse
// @Failure      500   {object}  dto.APIResponse
// @Router       /{code} [get]
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
	c.Redirect(http.StatusMovedPermanently, urlObj.LongURL)
}
