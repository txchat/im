#!/bin/bash
work_dir=$(
    cd "$(dirname "$0")" || exit
    pwd
)/../..

echo "${work_dir}"

protoc -I. -I"$GOPATH"/src --go_out=. --go_opt=paths=source_relative protocol/protocol.proto
protoc -I. -I"$GOPATH"/src --go_out=. --go-grpc_out=. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative comet/comet.proto
protoc -I. -I"$GOPATH"/src --go_out=. --go-grpc_out=. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative logic/logic.proto
