package operationlog

import (
	"github.com/gin-gonic/gin"

	"gobackend/pkg/core"
	metav1 "gobackend/pkg/meta/v1"
)

// Delete deletes an operation log record.
func (o *Controller) Delete(c *gin.Context) {
	if err := o.store.OperationLogs().Delete(
		c,
		c.Param("id"),
		metav1.DeleteOptions{Unscoped: true},
	); err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, nil)
}
