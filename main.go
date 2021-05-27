package main

import (
	"fmt"
	"os"

	"github.com/ren-kt/split_downloader/downloading"
)

func main() {
	fmt.Println("Download started")
	client := &downloading.Client{Download: downloading.NewDownload()}

	status := client.Run()
	switch status {
	case downloading.StatusOK:
		fmt.Println("Download complete")
	case downloading.StatusErr:
		fmt.Println("Download abnormal termination")
	}
	os.Exit(status)
}
