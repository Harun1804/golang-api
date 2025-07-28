package helpers

import (
	"github.com/gin-gonic/gin"
)

type SuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type ErrorResponse struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors"`
}

// Helper to send error response
func SendError(ctx *gin.Context, status int, msg string, err error) {
	ctx.JSON(status, ErrorResponse{
		Success: false,
		Message: msg,
		Errors:  TranslateErrorMessage(err),
	})
}

// Helper to send success response
func SendSuccess(ctx *gin.Context, status int, msg string, data interface{}) {
	ctx.JSON(status, SuccessResponse{
		Success: true,
		Message: msg,
		Data:    data,
	})
}