package posixsignal

import (
	"syscall"
	"testing"
	"time"

	"go-web-backend/pkg/shutdown"
)

type startShutdownFunc func(sm shutdown.Manager)

func (f startShutdownFunc) StartShutdown(sm shutdown.Manager) {
	f(sm)
}

func (f startShutdownFunc) ReportError(err error) {

}

func (f startShutdownFunc) AddShutdownCallback(shutdownCallback shutdown.Callback) {

}

func waitSig(t *testing.T, c <-chan int) {
	select {
	case <-c:

	case <-time.After(1 * time.Second):
		t.Error("Timeout waiting for StartShutdown.")
	}
}

func TestStartShutdownCalledOnDefaultSignals(t *testing.T) {
	c := make(chan int, 100)

	psm := NewPosixSignalManager()
	_ = psm.Start(startShutdownFunc(func(sm shutdown.Manager) {
		c <- 1
	}))

	time.Sleep(time.Millisecond)

	_ = syscall.Kill(syscall.Getpid(), syscall.SIGINT)

	waitSig(t, c)

	_ = psm.Start(startShutdownFunc(func(sm shutdown.Manager) {
		c <- 1
	}))

	time.Sleep(time.Millisecond)

	_ = syscall.Kill(syscall.Getpid(), syscall.SIGTERM)

	waitSig(t, c)
}

func TestStartShutdownCalledCustomSignal(t *testing.T) {
	c := make(chan int, 100)

	psm := NewPosixSignalManager(syscall.SIGHUP)
	_ = psm.Start(startShutdownFunc(func(sm shutdown.Manager) {
		c <- 1
	}))

	time.Sleep(time.Millisecond)

	_ = syscall.Kill(syscall.Getpid(), syscall.SIGHUP)

	waitSig(t, c)
}
