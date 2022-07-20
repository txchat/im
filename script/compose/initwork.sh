#!/bin/bash
# shellcheck disable=SC2034
# shellcheck disable=SC2048
work_dir=$(
    cd "$(dirname "$0")" || exit
    pwd
)/../..

created_volume=()
created_network=()

serviceList=$1
projectVersion=$2

# platform adaptation
HOST_OS=$(uname -s | tr '[:upper:]' '[:lower:]')
case ${HOST_OS} in
    "darwin")
        function dosed() {
            sed -r -i '' "$@"
        }
        ;;
    *)
        function dosed() {
            sed -r -i "$@"
        }
        ;;
esac

function network_create() {
    networkName=$1
    filterName=$(docker network ls | awk 'NR>1{ print $2 }' | grep -w "${networkName}")
    if [ "$filterName" == "" ]; then
        #不存在就创建
        created_network+=("$networkName")
        docker network create "$networkName"
    fi
}

function volume_create() {
    volumeName=$1
    filterName=$(docker volume ls | awk 'NR>1{ print $2 }' | grep -w "${volumeName}")
    if [ "$filterName" == "" ]; then
        #不存在就创建
        created_volume+=("$volumeName")
        docker volume create "$volumeName"
    fi
}

function initRedis() {
    for vname in ${created_volume[*]}; do
        if [ "${vname}" = "txchat-redis-config" ]; then
            docker container create --name dummy -v "txchat-redis-config":/root hello-world
            docker cp redis.conf dummy:/root/redis.conf
            docker rm dummy
        fi
    done
}

volumes=("txchat-zookeeper" "txchat-kafka" "txchat-redis-data" "txchat-redis-config" "txchat-redis-log" "txchat-etcd-data")
networks=("txchat-components" "txchat-service")

for sName in ${serviceList}; do
    volumes+=("txchat-${sName}-config")
    # 将「-」转为「_」并将小写转大写
    upperSName=$(echo "${sName//[-]/_}" | tr '[:lower:]' '[:upper:]')
    # 修改.env文件服务镜像版本号
    dosed "s/(${upperSName}_IMAGE=)\s*(.+)/\1${projectVersion}/" .env
done

for vname in ${volumes[*]}; do
    volume_create "${vname}"
done

for netname in ${networks[*]}; do
    network_create "${netname}"
done

initRedis
