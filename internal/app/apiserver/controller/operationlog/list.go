package operationlog

import (
	"github.com/gin-gonic/gin"

	"gobackend/pkg/core"
	"gobackend/pkg/errors"
	"gobackend/pkg/fields"
	metav1 "gobackend/pkg/meta/v1"

	"gobackend/internal/pkg/code"
)

// List operation logs.
func (o *Controller) List(c *gin.Context) {
	var r metav1.ListOptions

	if err := c.ShouldBindQuery(&r); err != nil {
		core.WriteResponse(c, errors.WithCode(code.ErrBind, err.Error()), nil)

		return
	}

	if _, err := fields.ParseSelector(r.FieldSelector); err != nil {
		core.WriteResponse(c, errors.WithCode(code.ErrFieldSelectorValidation, ""), nil)

		return
	}

	operationLogs, err := o.store.OperationLogs().List(c, r)

	if err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, operationLogs)
}
