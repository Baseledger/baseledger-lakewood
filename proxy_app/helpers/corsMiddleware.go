package helpers

import "github.com/gin-gonic/gin"

const authorizationHeader = "Authorization"
const defaultCorsAccessControlAllowOrigin = "*"
const defaultCorsAccessControlAllowCredentials = "true"
const defaultCorsAccessControlAllowHeaders = "Accept, Accept-Encoding, Authorization, Cache-Control, Content-Length, Content-Type, Origin, User-Agent, X-CSRF-Token, X-Requested-With"
const defaultCorsAccessControlAllowMethods = "GET, POST, PUT, DELETE, OPTIONS"
const defaultCorsAccessControlExposeHeaders = "X-Total-Results-Count"
const defaultResponseContentType = "application/json; charset=UTF-8"
const defaultResultsPerPage = 25

// CORSMiddleware is a working middlware for using CORS with gin
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", defaultCorsAccessControlAllowOrigin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", defaultCorsAccessControlAllowCredentials)
		c.Writer.Header().Set("Access-Control-Allow-Headers", defaultCorsAccessControlAllowHeaders)
		c.Writer.Header().Set("Access-Control-Allow-Methods", defaultCorsAccessControlAllowMethods)
		c.Writer.Header().Set("Access-Control-Expose-Headers", defaultCorsAccessControlExposeHeaders)

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
