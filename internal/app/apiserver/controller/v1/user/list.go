package user

import (
	"github.com/gin-gonic/gin"

	"gobackend/pkg/core"
	"gobackend/pkg/errors"
	"gobackend/pkg/fields"
	"gobackend/pkg/log"
	metav1 "gobackend/pkg/meta/v1"

	"gobackend/internal/pkg/code"
)

// List users.
func (u *Controller) List(c *gin.Context) {
	log.C(c).Debug("list domains function called")

	var r metav1.ListOptions
	if err := c.ShouldBindQuery(&r); err != nil {
		core.WriteResponse(c, errors.WithCode(code.ErrBind, err.Error()), nil)

		return
	}

	if _, err := fields.ParseSelector(r.FieldSelector); err != nil {
		core.WriteResponse(c, errors.WithCode(code.ErrFieldSelectorValidation, ""), nil)

		return
	}

	users, err := u.srv.Users().List(c, r)

	if err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, users)
}
