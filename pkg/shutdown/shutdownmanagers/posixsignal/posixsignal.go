package posixsignal

import (
	"os"
	"os/signal"
	"syscall"

	"go-web-demo/pkg/shutdown"
)

// Name defines shutdown manager name.
const Name = "PosixSignalManager"

// Manager implements ShutdownManager interface that is added
// to GracefulShutdown. Initialize with NewPosixSignalManager.
type Manager struct {
	signals []os.Signal
}

// NewPosixSignalManager initializes the PosixSignalManager.
// As arguments you can provide os.Signal-s to listen to, if none are given,
// it will default to SIGINT and SIGTERM.
func NewPosixSignalManager(sig ...os.Signal) *Manager {
	if len(sig) == 0 {
		sig = make([]os.Signal, 2)
		sig[0] = os.Interrupt
		sig[1] = syscall.SIGTERM
	}

	return &Manager{
		signals: sig,
	}
}

// GetName returns name of this ShutdownManager.
func (posixSignalManager *Manager) GetName() string {
	return Name
}

// Start starts listening for posix signals.
func (posixSignalManager *Manager) Start(gs shutdown.GSInterface) error {
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, posixSignalManager.signals...)

		// Block until a signal is received.
		<-c

		gs.StartShutdown(posixSignalManager)
	}()

	return nil
}

// ShutdownStart does nothing.
func (posixSignalManager *Manager) ShutdownStart() error {
	return nil
}

// ShutdownFinish exits the app with os.Exit(0).
func (posixSignalManager *Manager) ShutdownFinish() error {
	os.Exit(0)

	return nil
}
