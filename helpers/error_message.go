package helpers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// ParseError parses the error message and returns a map of field errors
func ParseError(err error) map[string]string {
	errors := make(map[string]string)

	if err, ok := err.(validator.ValidationErrors); ok {
		for _, e := range err {
			field := e.Field()
			tag := e.Tag()

			// Customize error messages as needed
			var message string
			switch tag {
			case "required":
				message = fmt.Sprintf("%s is required", field)
			case "min":
				message = fmt.Sprintf("%s must be at least %s", field, e.Param())
			case "oneof":
				message = fmt.Sprintf("%s only allows %s", field, e.Param())
			default:
				message = fmt.Sprintf("%s is not valid", field)
			}

			errors[field] = message
		}
	}

	return errors
}

// ErrorWithDataJSON returns a JSON response with error data
func HelperErrorWithDataJSON(c *gin.Context, errors map[string]string) {
	var errorList []map[string]string

	for field, message := range errors {
		errorItem := map[string]string{
			"field":   field,
			"message": message,
		}
		errorList = append(errorList, errorItem)
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"status":  "error",
		"message": "Validation failed",
		"data":    errorList,
	})
}
