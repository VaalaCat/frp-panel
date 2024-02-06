#!/usr/bin/env /bin/bash
bash ./codegen.sh
mkdir -p dist
rm -rf dist/*
cd www && pnpm install && pnpm build && cd ..
echo "Building frp-panel full windows binaries..."
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o dist/frp-panel-amd64.exe cmd/frpp/*.go
CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -o dist/frp-panel-arm64.exe cmd/frpp/*.go
echo "Building frp-panel full linux binaries..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dist/frp-panel-linux-amd64 cmd/frpp/*.go
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o dist/frp-panel-linux-arm64 cmd/frpp/*.go
echo "Building frp-panel full darwin binaries..."
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o dist/frp-panel-darwin-amd64 cmd/frpp/*.go
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o dist/frp-panel-darwin-arm64 cmd/frpp/*.go

echo "Building frp-panel client only windows binaries..."
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o dist/frp-panel-client-amd64.exe cmd/frppc/*.go
CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -o dist/frp-panel-client-arm64.exe cmd/frppc/*.go
echo "Building frp-panel client only linux binaries..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dist/frp-panel-client-linux-amd64 cmd/frppc/*.go
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o dist/frp-panel-client-linux-arm64 cmd/frppc/*.go
echo "Building frp-panel client only darwin binaries..."
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o dist/frp-panel-client-darwin-amd64 cmd/frpp/*.go
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o dist/frp-panel-client-darwin-arm64 cmd/frppc/*.go

echo "Build Done!"