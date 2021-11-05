package user

import (
	"github.com/gin-gonic/gin"

	"go-web-backend/pkg/core"
	"go-web-backend/pkg/errors"
	"go-web-backend/pkg/log"
	metav1 "go-web-backend/pkg/meta/v1"

	"go-web-backend/internal/pkg/code"
	v1 "go-web-backend/internal/pkg/entity/apiserver/v1"
)

// Update update a user info by the user identifier.
func (u *Controller) Update(c *gin.Context) {
	log.C(c).Info("update user function called.")

	var r v1.User

	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteResponse(c, errors.WithCode(code.ErrBind, err.Error()), nil)

		return
	}

	user, err := u.srv.Users().Get(c, c.Param("name"), metav1.GetOptions{})
	if err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	user.Nickname = r.Nickname
	user.Email = r.Email
	user.Phone = r.Phone
	user.Extend = r.Extend

	if errs := user.ValidateUpdate(); len(errs) != 0 {
		core.WriteResponse(c, errors.WithCode(code.ErrValidation, errs.ToAggregate().Error()), nil)

		return
	}

	// Save changed fields.
	if err := u.srv.Users().Update(c, user, metav1.UpdateOptions{}); err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, user)
}
