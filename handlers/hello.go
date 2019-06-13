package handlers

import (
	"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
)

func HelloHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	c.JSON(200, gin.H{
		"userID":   claims["id"],
		"text":     "Hello World.",
	})
}
