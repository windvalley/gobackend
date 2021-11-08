package user

import (
	srvv1 "gobackend/internal/app/apiserver/service/v1"
	"gobackend/internal/app/apiserver/store"
)

// Controller create a user handler used to handle request for user resource.
type Controller struct {
	srv srvv1.Service
}

// NewController creates a user handler.
func NewController(store store.Factory) *Controller {
	return &Controller{
		srv: srvv1.NewService(store),
	}
}
