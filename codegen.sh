#!/usr/bin/env /bin/bash

PROTOC_PATH=$(whereis protoc | awk '{print $2}')

cd idl && $PROTOC_PATH *.proto --go_out=. --go-grpc_out=. && cd ..
cd www && npx $PROTOC_PATH --ts_out ./lib/pb -I ../idl --proto_path ../idl/common.proto ../idl/common.proto ../idl/api*.proto && cd ..