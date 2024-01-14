#!/usr/bin/env /bin/bash
mkdir -p dist
GOOS=windows GOARCH=amd64 go build -o dist/frp-panel-amd64.exe cmd/*.go
GOOS=windows GOARCH=arm64 go build -o dist/frp-panel-arm64.exe cmd/*.go
GOOS=linux GOARCH=amd64 go build -o dist/frp-panel-linux-amd64 cmd/*.go
GOOS=linux GOARCH=arm64 go build -o dist/frp-panel-linux-arm64 cmd/*.go
GOOS=darwin GOARCH=amd64 go build -o dist/frp-panel-darwin-amd64 cmd/*.go
GOOS=darwin GOARCH=arm64 go build -o dist/frp-panel-darwin-arm64 cmd/*.go