package middlewares

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ValidateIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID, try again!"})
			c.Abort()
			return
		}
		c.Set("id", id)
		c.Next()
	}
}
