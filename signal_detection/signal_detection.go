package signal_detection

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var trapSignals = []os.Signal{
	syscall.SIGHUP,
	syscall.SIGINT,
	syscall.SIGTERM,
	syscall.SIGQUIT,
}

var OsExit = os.Exit

var DeleteTempFile = func() {}

func Listen(timeout time.Duration) (context.Context, func()) {
	bc := context.Background()
	ctx, cancel := context.WithTimeout(bc, timeout)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, trapSignals...)

	go func() {
		sig := <-sigCh
		fmt.Println("Got signal", sig)
		DeleteTempFile()
		cancel()
		OsExit(0)
	}()

	return ctx, cancel
}
