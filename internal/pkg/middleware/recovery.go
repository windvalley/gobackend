package middleware

import (
	"fmt"
	"runtime/debug"

	"github.com/fatih/color"
	"github.com/gin-gonic/gin"

	"gobackend/pkg/core"
	"gobackend/pkg/errors"

	"gobackend/internal/pkg/code"
)

// Recovery from panic.
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				debugStack := fmt.Sprintf("%s\n%s", err, debug.Stack())

				core.WriteResponse(c, errors.WithCode(code.ErrUnknown, debugStack), "")

				if gin.Mode() == gin.DebugMode {
					fmt.Printf(
						"%s\n%s",
						color.RedString(fmt.Sprintf("%s", err)),
						color.MagentaString(string(debug.Stack())),
					)
				}
			}
		}()

		c.Next()
	}
}
