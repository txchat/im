#!/bin/bash

# 获取编译参数
getFlags() {
  main_path="main"
  go_version=$(go version | awk '{ print $3 }')
  build_time=$(date "+%Y-%m-%d %H:%M:%S %Z")
  git_commit=$(git rev-parse --short=10 HEAD)
  flags="-X '${main_path}.goVersion=${go_version}' -X '${main_path}.buildTime=${build_time}' -X '${main_path}.gitCommit=${git_commit}' -X 'google.golang.org/protobuf/reflect/protoregistry.conflictPolicy=warn'"
  echo "${flags}"
}

# 设置目标打包环境
# 默认 linux amd64
initOS() {
  env_type="amd64"
  if [ -n "$1" ]; then
     env_type="$1"
  fi

  export GOOS=linux
  export GOARCH=${env_type}
  export GOLANG_PROTOBUF_REGISTRATION_CONFLICT=warn

  echo "linux_${env_type}"
}