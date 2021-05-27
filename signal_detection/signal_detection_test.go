package signal_detection_test

import (
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/ren-kt/split_downloader/signal_detection"
)

func TestListen(t *testing.T) {
	cases := []struct {
		name   string
		signal bool
	}{
		{
			name:   "execute Ctrl+C",
			signal: true,
		},
		{
			name:   "timeout",
			signal: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx, _ := signal_detection.Listen(1 * time.Millisecond)

			if c.signal {
				doneCh := make(chan int)
				signal_detection.OsExit = func(code int) { doneCh <- code }

				process, err := os.FindProcess(os.Getpid())
				if err != nil {
					t.Fatal(err)
				}

				err = process.Signal(syscall.SIGTERM)
				if err != nil {
					t.Fatal(err)
				}
			}

			select {
			case <-ctx.Done():
				return
			case <-time.After(2 * time.Millisecond):
				t.Error("Cancellation failure")
			}
		})
	}
}
