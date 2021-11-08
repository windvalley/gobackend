package user

import (
	"github.com/gin-gonic/gin"

	"gobackend/pkg/core"
	"gobackend/pkg/log"
	metav1 "gobackend/pkg/meta/v1"
)

// DeleteCollection batch delete users by multiple usernames.
// Only administrator can call this function.
func (u *Controller) DeleteCollection(c *gin.Context) {
	log.C(c).Info("batch delete user function called.")

	usernames := c.QueryArray("name")

	if err := u.srv.Users().DeleteCollection(c, usernames, metav1.DeleteOptions{}); err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, nil)
}
