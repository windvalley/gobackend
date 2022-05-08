package operationlog

import (
	"gobackend/internal/app/apiserver/store"
)

// Controller is the operationlog controller.
type Controller struct {
	store store.Factory
}

// NewController returns a operationlog controller.
func NewController(store store.Factory) *Controller {
	return &Controller{
		store: store,
	}
}
