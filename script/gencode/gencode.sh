#!/bin/bash
work_dir=$(
    cd "$(dirname "$0")" || exit
    pwd
)/../..

echo "${work_dir}"

service="comet logic"

# shellcheck disable=SC2154
for serverName in ${service}; do
    targetDir="app/${serverName}"
    goctl rpc protoc "${work_dir}/${targetDir}/${serverName}.proto" \
        --go_opt=module="github.com/txchat/im" --go-grpc_opt=module="github.com/txchat/im" \
        --proto_path="${work_dir}" --go_out="${work_dir}" --go-grpc_out="${work_dir}" --zrpc_out="${work_dir}/${targetDir}"
done
