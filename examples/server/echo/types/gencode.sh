#!/bin/bash
work_dir=$(
    cd "$(dirname "$0")" || exit
    pwd
)/../..

echo "${work_dir}"

protoc -I. -I"$GOPATH"/src --go_out=. --go_opt=paths=source_relative echo.proto
