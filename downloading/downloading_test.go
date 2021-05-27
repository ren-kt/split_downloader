package downloading_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/ren-kt/split_downloader/downloading"
	"github.com/ren-kt/split_downloader/option"
)

var cases = []struct {
	name        string
	url         string
	output      string
	parallelNum int
	timeout     time.Duration
	result      int
}{
	{
		name:        "no parallel",
		url:         "test.com",
		output:      "no_parallel.txt",
		parallelNum: 1,
		timeout:     1 * time.Second,
		result:      0,
	},
	{
		name:        "parallel number 2",
		url:         "test.com",
		output:      "parallel_number_2.txt",
		parallelNum: 2,
		timeout:     1 * time.Second,
		result:      0,
	},
	{
		name:        "parallel number 100",
		url:         "test.com",
		output:      "parallel_number_100.txt",
		parallelNum: 100,
		timeout:     1 * time.Second,
		result:      0,
	},
	{
		name:        "timeout error",
		url:         "test.com",
		output:      "parallel_number_100.txt",
		parallelNum: 100,
		timeout:     1 * time.Millisecond,
		result:      1,
	},
}

func TestRun(t *testing.T) {
	// 別goroutine上でリッスンする
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Accept-Ranges", "bytes")
		fmt.Fprintln(w, "test test test test test test test")
	})
	ts := httptest.NewServer(h)
	defer ts.Close()

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			downloader := &downloading.Download{
				Options: &option.Options{
					URL:         ts.URL,
					Output:      c.output,
					ParallelNum: c.parallelNum,
					Timeout:     c.timeout,
				},
			}
			client := &downloading.Client{Download: downloader}
			if result := client.Run(); result != c.result {
				t.Errorf("want %d, got %d\n", c.result, result)
			}
		})
		deleteFile(t, c.output)
	}
}

func TestMockRun(t *testing.T) {
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			dummyDownloader := &downloading.DummyDownload{
				Options: &option.Options{
					URL:         c.url,
					Output:      c.output,
					ParallelNum: c.parallelNum,
					Timeout:     c.timeout,
				},
			}
			client := &downloading.Client{Download: dummyDownloader}
			if result := client.Run(); result != c.result {
				t.Errorf("want %d, got %d\n", c.result, result)
			}
		})
		deleteFile(t, c.output)
	}
}

func deleteFile(t *testing.T, f string) error {
	t.Helper()
	err := os.Remove(f)
	if err != nil {
		return err
	}
	return nil
}
