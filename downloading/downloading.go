package downloading

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"

	"github.com/ren-kt/split_downloader/option"
	"github.com/ren-kt/split_downloader/signal_detection"
	"golang.org/x/sync/errgroup"
)

const (
	StatusOK int = iota
	StatusErr
)

type Doer interface {
	getContentLength(ctx context.Context) (int, error)
	download(ctx context.Context, contentLength int, dir string) error
	parallelDownload(ctx context.Context, n, min, max int, dir string, errCh chan error)
	merge(dir string) error
}

type Client struct {
	Download Doer
}

type Download struct {
	Options *option.Options
}

func NewDownload() *Download {
	var options option.Options
	options.Parse()
	return &Download{
		Options: &options,
	}
}

func (d *Client) Run() int {
	// moacのためreflect使用
	options := reflect.ValueOf(d.Download).Elem().Field(0).Interface().(*option.Options)

	// シグナル受信
	ctx, cancel := signal_detection.Listen(options.Timeout)
	defer cancel()

	// contentLengthを取得
	contentLength, err := d.Download.getContentLength(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	// ダウンロード先の作成
	dir, err := ioutil.TempDir("", "download")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	// 添付ファイルを削除
	deleteTempFile := func() { os.RemoveAll(dir) }
	signal_detection.DeleteTempFile = deleteTempFile
	defer deleteTempFile()

	// ダウンロード実行
	err = d.Download.download(ctx, contentLength, dir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	// マージ実行
	err = d.Download.merge(dir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	return 0
}

func (d *Download) getContentLength(ctx context.Context) (int, error) {
	req, err := http.NewRequestWithContext(ctx, "HEAD", d.Options.URL, nil)
	if err != nil {
		return 0, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}

	if res.Header.Get("Accept-Ranges") != "bytes" {
		return int(res.ContentLength), nil
	} else if int(res.ContentLength) == 0 {
		return 0, errors.New("ContentLength is 0")
	} else {
		return int(res.ContentLength), nil
	}
}

func (d *Download) download(ctx context.Context, contentLength int, dir string) error {
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

	// for n := d.Options.ParallelNum; 0 < n; n-- {
	// 	n := n
	// 	g.Go(func() error {
	// 		min = preMin
	// 		max = contentLength/n - 1
	// 		preMin = contentLength / n
	// 		return d.parallelDownload(ctx, n, min, max, dir)
	// 	})
	// }

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}

func (d *Download) parallelDownload(ctx context.Context, n, min, max int, dir string, errCh chan error) {
	req, err := http.NewRequestWithContext(ctx, "GET", d.Options.URL, nil)
	if err != nil {
		errCh <- err
		return
	}
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", min, max))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		errCh <- err
		return
	}
	defer res.Body.Close()

	file, err := os.Create(fmt.Sprintf("%v/%v-%v", dir, n, d.Options.Output))

	if err != nil {
		errCh <- err
		return
	}
	defer file.Close()

	_, err = io.Copy(file, res.Body)
	if err != nil {
		errCh <- err
		return
	}

	errCh <- nil
}

func (d *Download) merge(dir string) error {
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
