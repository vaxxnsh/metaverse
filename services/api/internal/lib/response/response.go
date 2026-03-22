package response

import "github.com/gin-gonic/gin"

type SuccessResponse[T any] struct {
	Success bool `json:"success"`
	Data    T    `json:"data"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

func SendSuccess[T any](c *gin.Context, data T, statusCode int) {
	c.JSON(statusCode, SuccessResponse[T]{
		Success: true,
		Data:    data,
	})
}

func SendError(c *gin.Context, statusCode int, code, message string, details interface{}) {
	c.JSON(statusCode, ErrorResponse{
		Success: false,
		Code:    code,
		Message: message,
		Details: details,
	})
}
