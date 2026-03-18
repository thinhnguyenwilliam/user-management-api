// user-management-api/internal/utils/errors.go
package utils

import "github.com/gin-gonic/gin"

type ErrorCode string

const (
	ErrCodeBadRequest   ErrorCode = "BAD_REQUEST"
	ErrCodeInvalidInput ErrorCode = "INVALID_INPUT"
	ErrCodeValidation   ErrorCode = "VALIDATION_ERROR"
	ErrCodeUnauthorized ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden    ErrorCode = "FORBIDDEN"
	ErrCodeNotFound     ErrorCode = "NOT_FOUND"
	ErrCodeConflict     ErrorCode = "CONFLICT"

	ErrCodeInternal        ErrorCode = "INTERNAL_ERROR"
	ErrCodeDatabase        ErrorCode = "DATABASE_ERROR"
	ErrCodeCache           ErrorCode = "CACHE_ERROR"
	ErrCodeExternalService ErrorCode = "EXTERNAL_SERVICE_ERROR"

	ErrCodeInvalidToken ErrorCode = "INVALID_TOKEN"
	ErrCodeTokenExpired ErrorCode = "TOKEN_EXPIRED"
)

func httpStatusFromCode(code ErrorCode) int {
	switch code {

	case ErrCodeBadRequest,
		ErrCodeInvalidInput,
		ErrCodeValidation:
		return 400

	case ErrCodeUnauthorized,
		ErrCodeInvalidToken,
		ErrCodeTokenExpired:
		return 401

	case ErrCodeForbidden:
		return 403

	case ErrCodeNotFound:
		return 404

	case ErrCodeConflict:
		return 409

	case ErrCodeDatabase,
		ErrCodeCache,
		ErrCodeExternalService,
		ErrCodeInternal:
		return 500

	default:
		return 500
	}
}

type AppError struct {
	Message string
	Code    ErrorCode
	Err     error
}

func (ae *AppError) Error() string {
	return ae.Message
}

func (ae *AppError) Unwrap() error {
	return ae.Err
}

func NewError(message string, code ErrorCode) error {
	return &AppError{
		Message: message,
		Code:    code,
	}
}

func WrapError(message string, code ErrorCode, err error) error {
	return &AppError{
		Message: message,
		Code:    code,
		Err:     err,
	}
}

func ResponseError(ctx *gin.Context, err error) {
	appErr, ok := err.(*AppError)

	if !ok {
		ctx.JSON(500, gin.H{
			"error": "internal server error",
			"code":  ErrCodeInternal,
		})
		return
	}

	response := gin.H{
		"error": appErr.Message,
		"code":  appErr.Code,
	}

	if appErr.Err != nil {
		response["detail"] = appErr.Err.Error()
	}

	ctx.JSON(httpStatusFromCode(appErr.Code), response)
}

type APIResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    any         `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

func ResponseSuccess(ctx *gin.Context, status int, message string, data ...any) {

	resp := APIResponse{
		Status:  "success",
		Message: message,
	}

	if len(data) > 0 && data[0] != nil {
		resp.Data = data[0]
	}

	ctx.JSON(status, resp)
}

func ResponseSuccessWithMeta(ctx *gin.Context, status int, message string, data any, meta any) {
	ctx.JSON(status, APIResponse{
		Status:  "success",
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}
