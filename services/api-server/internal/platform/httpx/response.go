package httpx

import "github.com/gin-gonic/gin"

type SuccessResponse struct {
	RequestID string `json:"request_id"`
	Data      any    `json:"data"`
}

type ErrorResponse struct {
	RequestID string      `json:"request_id"`
	Error     ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code    string         `json:"code"`
	Message string         `json:"message"`
	Details map[string]any `json:"details,omitempty"`
}

func JSON(c *gin.Context, status int, data any) {
	c.JSON(status, SuccessResponse{
		RequestID: RequestIDFromContext(c),
		Data:      data,
	})
}

func AbortWithError(c *gin.Context, err error) {
	appErr := AsAppError(err)
	c.AbortWithStatusJSON(appErr.HTTPStatus, ErrorResponse{
		RequestID: RequestIDFromContext(c),
		Error: ErrorDetail{
			Code:    appErr.Code,
			Message: appErr.Message,
			Details: appErr.Details,
		},
	})
}
