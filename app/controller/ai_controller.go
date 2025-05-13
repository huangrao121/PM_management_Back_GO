package controller

import (
	"github.com/gin-gonic/gin"
)

type AiController interface {
	GenerateTask(c *gin.Context)
}
