package shutdown

import (
	"sync"
)

// Callback is an interface you have to implement for callbacks.
// OnShutdown will be called when shutdown is requested. The parameter
// is the name of the ShutdownManager that requested shutdown.
type Callback interface {
	OnShutdown(string) error
}

// Func is a helper type, so you can easily provide anonymous functions
// as Callbacks.
type Func func(string) error

// OnShutdown defines the action needed to run when shutdown triggered.
func (f Func) OnShutdown(shutdownManager string) error {
	return f(shutdownManager)
}

// Manager is an interface implemnted by ShutdownManagers.
// GetName returns the name of ShutdownManager.
// ShutdownManagers start listening for shutdown requests in Start.
// When they call StartShutdown on GSInterface,
// first ShutdownStart() is called, then all ShutdownCallbacks are executed
// and once all ShutdownCallbacks return, ShutdownFinish is called.
type Manager interface {
	GetName() string
	Start(gs GSInterface) error
	ShutdownStart() error
	ShutdownFinish() error
}

// ErrorHandler is an interface you can pass to SetErrorHandler to
// handle asynchronous errors.
type ErrorHandler interface {
	OnError(err error)
}

// ErrorFunc is a helper type, so you can easily provide anonymous functions
// as ErrorHandlers.
type ErrorFunc func(err error)

// OnError defines the action needed to run when error occurred.
func (f ErrorFunc) OnError(err error) {
	f(err)
}

// GSInterface is an interface implemented by GracefulShutdown,
// that gets passed to ShutdownManager to call StartShutdown when shutdown
// is requested.
type GSInterface interface {
	StartShutdown(sm Manager)
	ReportError(err error)
	AddShutdownCallback(shutdownCallback Callback)
}

// GracefulShutdown is main struct that handles ShutdownCallbacks and
// ShutdownManagers. Initialize it with New.
type GracefulShutdown struct {
	callbacks    []Callback
	managers     []Manager
	errorHandler ErrorHandler
}

// New initializes GracefulShutdown.
func New() *GracefulShutdown {
	return &GracefulShutdown{
		callbacks: make([]Callback, 0, 10),
		managers:  make([]Manager, 0, 3),
	}
}

// Start calls Start on all added ShutdownManagers. The ShutdownManagers
// start to listen to shutdown requests. Returns an error if any ShutdownManagers
// return an error.
func (gs *GracefulShutdown) Start() error {
	for _, manager := range gs.managers {
		if err := manager.Start(gs); err != nil {
			return err
		}
	}

	return nil
}

// AddShutdownManager adds a ShutdownManager that will listen to shutdown requests.
func (gs *GracefulShutdown) AddShutdownManager(manager Manager) {
	gs.managers = append(gs.managers, manager)
}

// AddShutdownCallback adds a ShutdownCallback that will be called when
// shutdown is requested.
//
// You can provide anything that implements ShutdownCallback interface,
// or you can supply a function like this:
//	AddShutdownCallback(shutdown.ShutdownFunc(func() error {
//		// callback code
//		return nil
//	}))
func (gs *GracefulShutdown) AddShutdownCallback(shutdownCallback Callback) {
	gs.callbacks = append(gs.callbacks, shutdownCallback)
}

// SetErrorHandler sets an ErrorHandler that will be called when an error
// is encountered in ShutdownCallback or in ShutdownManager.
//
// You can provide anything that implements ErrorHandler interface,
// or you can supply a function like this:
//	SetErrorHandler(shutdown.ErrorFunc(func (err error) {
//		// handle error
//	}))
func (gs *GracefulShutdown) SetErrorHandler(errorHandler ErrorHandler) {
	gs.errorHandler = errorHandler
}

// StartShutdown is called from a ShutdownManager and will initiate shutdown.
// first call ShutdownStart on Shutdownmanager,
// call all ShutdownCallbacks, wait for callbacks to finish and
// call ShutdownFinish on ShutdownManager.
func (gs *GracefulShutdown) StartShutdown(sm Manager) {
	gs.ReportError(sm.ShutdownStart())

	var wg sync.WaitGroup
	for _, shutdownCallback := range gs.callbacks {
		wg.Add(1)
		go func(shutdownCallback Callback) {
			defer wg.Done()

			gs.ReportError(shutdownCallback.OnShutdown(sm.GetName()))
		}(shutdownCallback)
	}

	wg.Wait()

	gs.ReportError(sm.ShutdownFinish())
}

// ReportError is a function that can be used to report errors to
// ErrorHandler. It is used in ShutdownManagers.
func (gs *GracefulShutdown) ReportError(err error) {
	if err != nil && gs.errorHandler != nil {
		gs.errorHandler.OnError(err)
	}
}
