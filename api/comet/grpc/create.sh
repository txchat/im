#!/bin/sh

protoc -I. -I$GOPATH/src \
    --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    *.proto
