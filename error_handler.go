package go_task

import (
	"github.com/gin-gonic/gin"
	"log"
)

type ErrorHandler interface {
	HandleError(c *gin.Context, statusCode int, data gin.H)
}

type ErrorHandlerImpl struct{}

func (e *ErrorHandlerImpl) HandleError(c *gin.Context, statusCode int, data gin.H) {
	log.Printf("Error with status code %d: %v", statusCode, data)
	c.JSON(statusCode, data)
}
