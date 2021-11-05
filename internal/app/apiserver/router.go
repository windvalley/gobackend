package apiserver

import (
	"github.com/gin-gonic/gin"

	"go-web-demo/pkg/core"
	"go-web-demo/pkg/errors"
	"go-web-demo/pkg/log"

	"go-web-demo/internal/app/apiserver/controller/v1/user"
	"go-web-demo/internal/app/apiserver/store/mysql"
	"go-web-demo/internal/pkg/code"

	// Custom gin validators.
	_ "go-web-demo/internal/pkg/validator"
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
