package middleware

import (
	"github.com/gin-gonic/gin"
	gindump "github.com/tpkeeper/gin-dump"

	"gobackend/pkg/log"
)

// Dump header/body of request and response, very helpful for debugging your applications.
// NOTE: Do not use in production for performance consideration.
func Dump() gin.HandlerFunc {
	return gindump.DumpWithOptions(true, true, true, true, false, func(dumpStr string) {
		log.Info(dumpStr)
	})
}
