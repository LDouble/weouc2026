package httpx

import (
	"errors"
	"net/http"
)

const (
	CodeBadRequest   = "COMMON_BAD_REQUEST"
	CodeUnauthorized = "AUTH_UNAUTHORIZED"
	CodeForbidden    = "AUTH_FORBIDDEN"
	CodeNotFound     = "COMMON_NOT_FOUND"
	CodeUnavailable  = "COMMON_SERVICE_UNAVAILABLE"
	CodeInternal     = "COMMON_INTERNAL_ERROR"
)

type AppError struct {
	HTTPStatus int
	Code       string
	Message    string
	Details    map[string]any
	Err        error
}

func (e *AppError) Error() string {
	if e == nil {
		return ""
	}

	return e.Message
}

func (e *AppError) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Err
}

func BadRequest(message string, details map[string]any) *AppError {
	return &AppError{
		HTTPStatus: http.StatusBadRequest,
		Code:       CodeBadRequest,
		Message:    message,
		Details:    details,
	}
}

func Unauthorized(message string) *AppError {
	if message == "" {
		message = "需要登录后访问"
	}

	return &AppError{
		HTTPStatus: http.StatusUnauthorized,
		Code:       CodeUnauthorized,
		Message:    message,
	}
}

func Forbidden(message string, details map[string]any) *AppError {
	if message == "" {
		message = "当前账号无权访问该资源"
	}

	return &AppError{
		HTTPStatus: http.StatusForbidden,
		Code:       CodeForbidden,
		Message:    message,
		Details:    details,
	}
}

func NotFound(message string, details map[string]any) *AppError {
	if message == "" {
		message = "请求的资源不存在"
	}

	return &AppError{
		HTTPStatus: http.StatusNotFound,
		Code:       CodeNotFound,
		Message:    message,
		Details:    details,
	}
}

func Unavailable(message string) *AppError {
	if message == "" {
		message = "服务暂不可用"
	}

	return &AppError{
		HTTPStatus: http.StatusServiceUnavailable,
		Code:       CodeUnavailable,
		Message:    message,
	}
}

func Internal(message string, err error) *AppError {
	if message == "" {
		message = "服务器内部错误"
	}

	return &AppError{
		HTTPStatus: http.StatusInternalServerError,
		Code:       CodeInternal,
		Message:    message,
		Err:        err,
	}
}

func AsAppError(err error) *AppError {
	if err == nil {
		return nil
	}

	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}

	return Internal("", err)
}
