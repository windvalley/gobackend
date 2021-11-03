package user

import (
	"github.com/gin-gonic/gin"

	"go-web-demo/pkg/core"
	"go-web-demo/pkg/log"
	metav1 "go-web-demo/pkg/meta/v1"
)

// Delete delete an user by the user identifier.
// Only administrator can call this function.
func (u *Controller) Delete(c *gin.Context) {
	log.C(c).Info("delete user function called.")

	if err := u.srv.Users().Delete(c, c.Param("name"), metav1.DeleteOptions{Unscoped: true}); err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, nil)
}
