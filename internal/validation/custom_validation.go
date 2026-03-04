// user-management-api/internal/validation/custom_validation.go
package validation

import (
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

var (
	slugRegex   = regexp.MustCompile(`^[a-z0-9]+(?:[-.][a-z0-9]+)*$`)
	searchRegex = regexp.MustCompile(`^[a-zA-Z0-9\s]+$`)
)

func RegisterCustomValidation(v *validator.Validate) error {
	validations := map[string]validator.Func{
		"slug":     validateSlug,
		"search":   validateSearch,
		"min_int":  validateMinInt,
		"max_int":  validateMaxInt,
		"file_ext": validateFileExt,
	}

	for tag, fn := range validations {
		if err := v.RegisterValidation(tag, fn); err != nil {
			return err
		}
	}

	return nil
}

func validateSlug(fl validator.FieldLevel) bool {
	return slugRegex.MatchString(fl.Field().String())
}

func validateSearch(fl validator.FieldLevel) bool {
	return searchRegex.MatchString(fl.Field().String())
}

func validateMinInt(fl validator.FieldLevel) bool {
	minVal, err := strconv.ParseInt(fl.Param(), 10, 64)
	if err != nil {
		return false
	}
	return fl.Field().Int() >= minVal
}

func validateMaxInt(fl validator.FieldLevel) bool {
	maxVal, err := strconv.ParseInt(fl.Param(), 10, 64)
	if err != nil {
		return false
	}
	return fl.Field().Int() <= maxVal
}
func validateFileExt(fl validator.FieldLevel) bool {
	filename := fl.Field().String()
	allowedStr := fl.Param()
	if allowedStr == "" {
		return false
	}

	ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(filename)), ".")
	allowed := make(map[string]struct{})

	for _, a := range strings.Fields(allowedStr) {
		allowed[strings.ToLower(a)] = struct{}{}
	}

	_, exists := allowed[ext]
	return exists
}
