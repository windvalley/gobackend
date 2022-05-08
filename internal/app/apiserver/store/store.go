package store

var client Factory

// Factory is the interface of store client.
type Factory interface {
	Users() UserStore
	OperationLogs() OperationLogStore
	Close() error
}

// Client return the store client instance.
func Client() Factory {
	return client
}

// SetClient set the store client.
func SetClient(factory Factory) {
	client = factory
}
