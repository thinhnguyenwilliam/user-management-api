// user-management-api/internal/validation/validation.go
package validation

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/thinhnguyenwilliam/user-management-api/internal/utils"
)

func InitValidation() error {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return fmt.Errorf("failed to get validator engine")
	}

	return RegisterCustomValidation(v)
}

func HandleValidationErrors(err error) gin.H {
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return gin.H{"error": "Yêu cầu không hợp lệ: " + err.Error()}
	}

	errors := make(map[string]string)

	for _, e := range validationErrors {
		fieldPath := buildFieldPath(e.Namespace())
		errors[fieldPath] = translateError(fieldPath, e)
	}

	return gin.H{"error": errors}
}

func buildFieldPath(namespace string) string {
	parts := strings.Split(namespace, ".")

	if len(parts) > 1 {
		parts = parts[1:] // remove root struct
	}

	for i, part := range parts {
		if idx := strings.Index(part, "["); idx != -1 {
			base := utils.CamelToSnake(part[:idx])
			parts[i] = base + part[idx:]
			continue
		}
		parts[i] = utils.CamelToSnake(part)
	}

	return strings.Join(parts, ".")
}

func translateError(field string, e validator.FieldError) string {
	switch e.Tag() {

	case "gt":
		return fmt.Sprintf("%s phải lớn hơn %s", field, e.Param())

	case "lt":
		return fmt.Sprintf("%s phải nhỏ hơn %s", field, e.Param())

	case "gte":
		return fmt.Sprintf("%s phải lớn hơn hoặc bằng %s", field, e.Param())

	case "lte":
		return fmt.Sprintf("%s phải nhỏ hơn hoặc bằng %s", field, e.Param())

	case "uuid":
		return fmt.Sprintf("%s phải là UUID hợp lệ", field)

	case "slug":
		return fmt.Sprintf("%s chỉ được chứa chữ thường, số, dấu gạch ngang hoặc dấu chấm", field)

	case "min":
		return fmt.Sprintf("%s phải nhiều hơn %s ký tự", field, e.Param())

	case "max":
		return fmt.Sprintf("%s phải ít hơn %s ký tự", field, e.Param())

	case "min_int":
		return fmt.Sprintf("%s phải có giá trị lớn hơn %s", field, e.Param())

	case "max_int":
		return fmt.Sprintf("%s phải có giá trị bé hơn %s", field, e.Param())

	case "oneof":
		values := strings.Join(strings.Split(e.Param(), " "), ",")
		return fmt.Sprintf("%s phải là một trong các giá trị: %s", field, values)

	case "required":
		return fmt.Sprintf("%s là bắt buộc", field)

	case "search":
		return fmt.Sprintf("%s chỉ được chứa chữ thường, in hoa, số và khoảng trắng", field)

	case "email":
		return fmt.Sprintf("%s phải đúng định dạng email", field)

	case "datetime":
		return fmt.Sprintf("%s phải theo định dạng YYYY-MM-DD", field)

	case "file_ext":
		values := strings.Join(strings.Split(e.Param(), " "), ",")
		return fmt.Sprintf("%s chỉ cho phép extension: %s", field, values)

	default:
		return fmt.Sprintf("%s không hợp lệ", field)
	}
}
