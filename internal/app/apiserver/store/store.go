package store

var client Factory

// Factory ...
type Factory interface {
	Users() UserStore
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
