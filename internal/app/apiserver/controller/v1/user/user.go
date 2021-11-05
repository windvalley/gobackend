package user

import (
	srvv1 "go-web-backend/internal/app/apiserver/service/v1"
	"go-web-backend/internal/app/apiserver/store"
)

// Controller create a user handler used to handle request for user resource.
type Controller struct {
	srv srvv1.Service
}

// NewUserController creates a user handler.
func NewUserController(store store.Factory) *Controller {
	return &Controller{
		srv: srvv1.NewService(store),
	}
}
