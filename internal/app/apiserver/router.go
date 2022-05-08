package apiserver

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"gobackend/pkg/core"
	"gobackend/pkg/errors"
	"gobackend/pkg/log"

	"gobackend/internal/app/apiserver/controller/operationlog"
	"gobackend/internal/app/apiserver/controller/v1/user"
	"gobackend/internal/app/apiserver/store/mysql"
	"gobackend/internal/pkg/code"
	"gobackend/internal/pkg/middleware"

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
		core.WriteResponse(c, errors.WithCode(code.ErrPageNotFound, "URL path not found"), nil)
	})

	storeIns := mysql.GetMysqlFactory()

	log.Infof("get mysql factory instance: %v", storeIns)

	// Operation logging.
	if viper.GetBool("feature.operation-logging") {
		g.Use(middleware.OperationLog(storeIns))

		ol := g.Group("/operation-logs")
		{
			olController := operationlog.NewController(storeIns)

			ol.GET("", olController.List)
			ol.DELETE(":id", olController.Delete)
		}
	}

	v1 := g.Group("/v1")
	{
		userv1 := v1.Group("/users")
		{
			userController := user.NewController(storeIns)

			userv1.POST("", userController.Create)
			userv1.GET(":name", userController.Get)
			userv1.GET("", userController.List)
			userv1.PUT(":name", userController.Update)
			userv1.DELETE(":name", userController.Delete)
			userv1.DELETE("", userController.DeleteCollection)
		}
	}

	return g
}
