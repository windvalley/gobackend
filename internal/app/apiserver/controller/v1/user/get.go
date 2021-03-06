package user

import (
	"github.com/gin-gonic/gin"

	"gobackend/pkg/core"
	"gobackend/pkg/log"
	metav1 "gobackend/pkg/meta/v1"
)

// Get get an user by the user identifier.
func (u *Controller) Get(c *gin.Context) {
	log.C(c).Debug("user Get function is called")

	user, err := u.srv.Users().Get(c, c.Param("name"), metav1.GetOptions{})
	if err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, user)
}
