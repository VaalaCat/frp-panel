#!/usr/bin/env /bin/bash
bash ./codegen.sh
mkdir -p dist
rm -rf dist/*
cd www && pnpm install && pnpm build && cd ..
echo "Building frp-panel windows binaries..."
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o dist/frp-panel-amd64.exe cmd/*.go
CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -o dist/frp-panel-arm64.exe cmd/*.go
echo "Building frp-panel linux binaries..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dist/frp-panel-linux-amd64 cmd/*.go
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o dist/frp-panel-linux-arm64 cmd/*.go
echo "Building frp-panel darwin binaries..."
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o dist/frp-panel-darwin-amd64 cmd/*.go
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o dist/frp-panel-darwin-arm64 cmd/*.go

echo "Build Done!"