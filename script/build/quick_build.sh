#!/bin/bash
file_path=$(
  cd "$(dirname "$0")" || exit
  pwd
)/../..
# shellcheck source=./util.sh
source "${file_path}"/script/build/util.sh

quickBuildService() {
  file_path="$1"
  pkg_name="$2"
  service_name="$3"
  flags=$(getFlags)

  cd "${file_path}/${service_name}/cmd" || exit
  echo "┌ start building ${service_name} service"
  go build -ldflags "${flags}" -o "${file_path}/${pkg_name}/${service_name}" || exit
  echo "└ building ${service_name} service success"
}

mkDir() {
  file_path="$1"
  pkg_path="$2"

  rm -rf ${file_path}/${pkg_path}
  mkdir -pv "${file_path}/${pkg_path}"
}

env_type="$1"
initOS ${env_type}
os=$(initOS ${env_type})

api_version=""
project_name="im"
now_time=$(date "+%Y_%m_%d")
pkg_path="${project_name}_bin_${now_time}_${os}"

mkDir "${file_path}" "${pkg_path}"

quickBuildService "${file_path}" "${pkg_path}" "comet"
quickBuildService "${file_path}" "${pkg_path}" "logic"


