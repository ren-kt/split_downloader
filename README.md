# 分割ダウンローダ

## 仕様
- 分割ダウンロードを行う

## ダウンロードコマンドの例

```shell
$ docker-compose build
$ docker-compose run app /bin/bash
$ go build
$ ./split_downloader -u https://ajax.googleapis.com/ajax/libs/jquery/3.3.1/jquery.min.js -o download.js -p 2 -t 10s
```

## オプション

| オプション | 内容 | デフォルト |
| - | - | - |
| -u | 対象URL | https://ajax.googleapis.com/ajax/libs/jquery/3.3.1/jquery.min.js (ポートフォリオ用) |
| -o | 出力ファイル | download.js(ポートフォリオ用) |
| -p | 並列処理数 | 2 |
| -t | タイムアウトするまでの時間 | 10s |



 ## ディレクトリ構造

```
.
├── README.md
├── docker
│   └── go
│       └── Dockerfile
├── docker-compose.yml
├── download.js
├── downloading
│   ├── _downloading.go
│   ├── downloading.go
│   ├── downloading_test.go
│   └── export_test.go
├── go.mod
├── go.sum
├── main.go
├── master.js
├── option
│   └── option.go
└── signal_detection
    ├── signal_detection.go
    └── signal_detection_test.go
```