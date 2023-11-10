package http_util

import (
	"fmt"
	"gin-wire-template/utils/log_util"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
)

// HTTPSuccess Struct
type HTTPSuccess struct {
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"success message if any"`
}

// HTTPError Struct
type HTTPError struct {
	Status  string `json:"status" example:"fail"`
	Code    int    `json:"code" example:"500"`
	Message string `json:"message" example:"something went wrong"`
}

func RenderSuccessResponse(c *gin.Context, obj any, statusCode ...int) {
	code := http.StatusOK

	if len(statusCode) > 0 {
		code = statusCode[0]
	}

	if obj != nil {
		if reflect.TypeOf(obj).Kind() == reflect.String {
			c.JSON(code, HTTPSuccess{
				Status:  "success",
				Message: obj.(string),
			})
			return
		}

		c.JSON(code, obj)
		return
	}

	//c.IndentedJSON
	c.JSON(code, HTTPSuccess{
		Status:  "success",
		Message: "",
	})
	return
}

func RenderErrorResponse(c *gin.Context, obj any, statusCode ...int) {
	code := http.StatusInternalServerError

	if len(statusCode) > 0 {
		code = statusCode[0]
	}

	if obj != nil {
		if reflect.TypeOf(obj).Kind() == reflect.String {
			log_util.Logger.Warn(fmt.Sprintf("Error response: %v", obj))
			c.JSON(code, HTTPError{
				Status:  "fail",
				Code:    code,
				Message: obj.(string),
			})
			return
		}

		c.JSON(code, obj)
		return
	}

	c.JSON(code, HTTPError{
		Status:  "fail",
		Code:    code,
		Message: "",
	})
	return
}
