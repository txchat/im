#!/bin/sh
path=$(
    cd "$(dirname "$0")" || exit
    pwd
)

# machine x
# ARCH=$(uname -m)
OS=$(uname -s | tr '[:upper:]' '[:lower:]')

case ${OS} in
    "darwin")
        export PATH="${path}/osx_x86_64":"$PATH"
        ;;
    "linux")
        export PATH="${path}/linux_x86_64":"$PATH"
        ;;
esac
