package option

import (
	"flag"
	"time"
)

type Options struct {
	URL         string
	Output      string
	ParallelNum int
	Timeout     time.Duration
}

func (options *Options) Parse() {
	// ポートフォリオ用
	flag.StringVar(&options.URL, "u", "https://ajax.googleapis.com/ajax/libs/jquery/3.3.1/jquery.min.js", "Download destination url")
	flag.StringVar(&options.Output, "o", "download.js", "Download file name")
	flag.IntVar(&options.ParallelNum, "p", 2, "Paralle number")
	flag.DurationVar(&options.Timeout, "t", 10*time.Second, "Timeout")
	flag.Parse()
}
