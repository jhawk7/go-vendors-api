package handlers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func ErrorHandler(c *gin.Context, err error, status int, fatal bool) {
	if err != nil {
		log.Error(fmt.Errorf("error: %v", err.Error()))

		if fatal {
			panic(err)
		}

		if status != 0 {
			c.AbortWithStatusJSON(status, gin.H{
				"error": err.Error(),
			})
		}
	}
}

func LogInfo(info string) {
	log.Info(info)
}
