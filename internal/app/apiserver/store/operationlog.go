package store

import (
	"context"

	metav1 "gobackend/pkg/meta/v1"

	"gobackend/internal/pkg/entity/apiserver/operationlog"
)

// OperationLogStore is an interface for storing operation logs.
type OperationLogStore interface {
	Create(
		ctx context.Context,
		operationLog *operationlog.OperationLog,
		opts metav1.CreateOptions,
	) error

	List(
		ctx context.Context,
		opts metav1.ListOptions,
	) (*operationlog.List, error)

	Delete(
		ctx context.Context,
		id string,
		opts metav1.DeleteOptions,
	) error
}
