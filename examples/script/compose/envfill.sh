#!/bin/bash
# shellcheck disable=SC2034
work_dir=$(
    cd "$(dirname "$0")" || exit
    pwd
)/../..

PROJECT_VERSION=$1

function randomPassword() {
    return
}

if [ -f .env ]; then
    echo ".env file existed"
else
    # shellcheck disable=SC1091
    source env_tmpl
    #source key 2>/dev/null || randomPassword

    eval "cat<<EOF > .env
$(<env_tmpl)
EOF"
fi
