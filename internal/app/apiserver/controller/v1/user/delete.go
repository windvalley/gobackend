package user

import (
	"github.com/gin-gonic/gin"

	"go-web-backend/pkg/core"
	"go-web-backend/pkg/log"
	metav1 "go-web-backend/pkg/meta/v1"
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
