package metahttp

import (
	"github.com/gin-gonic/gin"
)

var GetAllowedOrigins func() []string

func isAllowedOrigin(origin string) bool {
	if GetAllowedOrigins == nil {
		return false
	}
	allowedOrigins := GetAllowedOrigins()
	if len(allowedOrigins) == 0 {
		return false
	}
	for _, o := range allowedOrigins {
		if o == origin {
			return true
		}
	}
	return false
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if isAllowedOrigin(origin) {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
			c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, X-Token")
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
