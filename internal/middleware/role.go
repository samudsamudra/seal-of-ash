package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		r, ok := c.Get("role")
		if !ok || r != role {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "you are not bound to perform this ritual",
			})
			return
		}
		c.Next()
	}
}
