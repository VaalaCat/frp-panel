#!/usr/bin/env /bin/bash

cd idl && protoc *.proto --go_out=. --go-grpc_out=. && cd ..
cd www && npx protoc --ts_out ./lib/pb -I ../idl --proto_path ../idl/common.proto ../idl/common.proto ../idl/api*.proto && cd ..