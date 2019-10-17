#!/bin/bash -e

go test

if [[ -d releases ]]; then
  rm -rf releases
fi

mkdir releases

GOOS=linux GOARCH=amd64 go build -o releases/bosh-ip-util-linux-amd64 github.com/rahulkj/bosh-ip-util

GOOS=darwin GOARCH=amd64 go build -o releases/bosh-ip-util-darwin-amd64 github.com/rahulkj/bosh-ip-util

GOOS=windows GOARCH=386 go build -o releases/bosh-ip-util-windows-amd64.exe github.com/rahulkj/bosh-ip-util
