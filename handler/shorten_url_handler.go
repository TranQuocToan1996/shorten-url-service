package handler

import (
	"github.com/gin-gonic/gin"
)

type ShortenURLHandler struct{}

func NewFactorialHandler() *ShortenURLHandler {
	return &ShortenURLHandler{}
}

func (h *ShortenURLHandler) RequestCreateShortURL(c *gin.Context) {
	panic("not implemented")
	// var req dto.FactorialMessage

	// if err := c.ShouldBindJSON(&req); err != nil {
	// 	sendErrorResponse(c, http.StatusBadRequest, "fail", err.Error())
	// 	return
	// }

	// // Validate number
	// _, err := h.factorialService.ValidateNumber(fmt.Sprintf("%d", req.Number))
	// if err != nil {
	// 	sendErrorResponse(c, http.StatusBadRequest, "fail", err.Error())
	// 	return
	// }

	// msg := dto.FactorialMessage{Number: req.Number}
	// err = h.producer.Publish(c.Request.Context(), h.queueName, msg.Bytes())
	// if err != nil {
	// 	log.Printf("Error publishing message: %v", err)
	// 	sendErrorResponse(c, http.StatusInternalServerError, "fail", "Failed to submit calculation")
	// 	return
	// }

	// // Return calculating status
	// sendAPIResponse(c, http.StatusOK, "ok", "submitted", dto.CalculateResponseData{
	// 	Number:  req.Number,
	// 	Message: "submitted",
	// })
}
