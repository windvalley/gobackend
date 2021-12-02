package user

import (
	"github.com/gin-gonic/gin"

	"gobackend/pkg/core"
	"gobackend/pkg/log"
	metav1 "gobackend/pkg/meta/v1"
)

// Delete delete an user by the user identifier.
// Only administrator can call this function.
func (u *Controller) Delete(c *gin.Context) {
	log.C(c).Debug("delete user function called")

	if err := u.srv.Users().Delete(c, c.Param("name"), metav1.DeleteOptions{Unscoped: true}); err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, nil)
}
