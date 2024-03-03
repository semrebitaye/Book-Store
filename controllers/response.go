package controllers

import (
	"book-store/models"

	"github.com/gin-gonic/gin"
)

func SuccessResponse(c *gin.Context, status int, data any, meta *models.Metadata) {

	response := models.Response{
		Meta: meta,
		Data: data,
		Ok:   true,
	}
	c.JSON(status, response)
}

func ErrorResponse(c *gin.Context, status int, error *models.Error) {
	response := models.Response{
		Err: error,
		Ok:  false,
	}

	c.JSON(status, response)
}
