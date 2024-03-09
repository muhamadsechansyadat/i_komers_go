package helpers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type SuccessResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Error   bool        `json:"error"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Error   bool   `json:"error"`
}

type ErrorDataResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Error   bool        `json:"error"`
	Data    interface{} `json:"data,omitempty"`
}

func SuccessJSON(c *gin.Context, message string, data interface{}) {
	response := SuccessResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	}
	c.JSON(http.StatusOK, response)
}

func ErrorJSON(c *gin.Context, message string) {
	response := ErrorResponse{
		Status:  "error",
		Message: message,
	}
	c.JSON(http.StatusBadRequest, response)
}

func ErrorBadRequestJSON(c *gin.Context, message string) {
	response := ErrorResponse{
		Status:  "error",
		Message: message,
	}
	c.JSON(http.StatusBadRequest, response)
}

func SuccessStringJSON(c *gin.Context, message string) {
	response := ErrorResponse{
		Status:  "success",
		Message: message,
	}
	c.JSON(http.StatusBadRequest, response)
}

func ErrorWithDataJSON(c *gin.Context, message string, data interface{}) {
	response := ErrorDataResponse{
		Status:  "error",
		Message: message,
		Data:    data,
	}
	c.JSON(http.StatusBadRequest, response)
}
