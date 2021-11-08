package apiserver

import (
	"github.com/gin-gonic/gin"

	"gobackend/pkg/core"
	"gobackend/pkg/errors"
	"gobackend/pkg/log"

	"gobackend/internal/app/apiserver/controller/v1/user"
	"gobackend/internal/app/apiserver/store/mysql"
	"gobackend/internal/pkg/code"

	// Custom gin validators.
	_ "gobackend/internal/pkg/validator"
)

func initRouter(g *gin.Engine) {
	installMiddleware(g)
	installController(g)
}

func installMiddleware(g *gin.Engine) {
}

func installController(g *gin.Engine) *gin.Engine {
	g.NoRoute(func(c *gin.Context) {
		core.WriteResponse(c, errors.WithCode(code.ErrPageNotFound, "Page not found."), nil)
	})

	storeIns := mysql.GetMysqlFactory()

	log.Infof("get mysql factory instance: %v", storeIns)

	v1 := g.Group("/v1")
	{
		userv1 := v1.Group("/users")
		{
			userController := user.NewUserController(storeIns)

			userv1.POST("", userController.Create)
			userv1.DELETE("", userController.DeleteCollection)
			userv1.DELETE(":name", userController.Delete)
			userv1.PUT(":name", userController.Update)
			userv1.GET(":name", userController.Get)
		}
	}

	return g
}
