#!/usr/bin/env /bin/bash

cd idl && /opt/homebrew/Cellar/protobuf/25.1/bin/protoc *.proto --go_out=. --go-grpc_out=. && cd ..
cd www && npx /opt/homebrew/Cellar/protobuf/25.1/bin/protoc --ts_out ./lib/pb -I ../idl --proto_path ../idl/common.proto ../idl/common.proto ../idl/api*.proto && cd ..