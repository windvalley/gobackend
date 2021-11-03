package middleware

import (
	"github.com/gin-gonic/gin"

	"go-web-demo/pkg/log"
)

// UsernameKey defines the key in gin context which represents the owner of the secret.
const UsernameKey = "username"

// Context is a middleware that injects common prefix fields to gin.Context.
func Context() gin.HandlerFunc {
	return func(c *gin.Context) {
		l := log.WithFields(
			log.String("x-request-id", c.GetString(XRequestIDKey)),
			log.String("username", c.GetString(UsernameKey)),
		)

		c.Set(log.ContextLoggerName, l)

		c.Next()
	}
}
