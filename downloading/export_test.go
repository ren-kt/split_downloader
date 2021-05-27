package downloading

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ren-kt/split_downloader/option"
	"golang.org/x/sync/errgroup"
)

type DummyDownload struct {
	Options *option.Options
	Doer
}

func (d *DummyDownload) getContentLength(ctx context.Context) (int, error) {
	return 35, nil
}

func (d *DummyDownload) download(ctx context.Context, contentLength int, dir string) error {
	var preMin, min, max int
	errCh := make(chan error)

	g, ctx := errgroup.WithContext(ctx)

	for n := d.Options.ParallelNum; 0 < n; n-- {
		n := n
		min = preMin
		max = contentLength/n - 1
		preMin = contentLength / n
		go d.parallelDownload(ctx, n, min, max, dir, errCh)
	}

	for n := d.Options.ParallelNum; 0 < n; n-- {
		g.Go(func() error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case err := <-errCh:
				return err
			}
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}

func (d *DummyDownload) parallelDownload(ctx context.Context, n, min, max int, dir string, errCh chan error) {
	file, err := os.Create(fmt.Sprintf("%v/%v-%v", dir, n, d.Options.Output))

	if err != nil {
		errCh <- err
		return
	}
	defer file.Close()

	_, err = io.Copy(file, strings.NewReader("test"))
	if err != nil {
		errCh <- err
		return
	}

	errCh <- nil
}

func (d *DummyDownload) merge(dir string) error {
	file, err := os.Create(d.Options.Output)
	if err != nil {
		return err
	}
	defer file.Close()

	for n := d.Options.ParallelNum; 0 < n; n-- {
		src, err := os.Open(fmt.Sprintf("%v/%v-%v", dir, n, d.Options.Output))
		if err != nil {
			return err
		}
		_, err = io.Copy(file, src)
		if err != nil {
			return err
		}
	}
	return nil
}
