#!/bin/bash
work_dir=$(
  cd "$(dirname "$0")" || exit
  pwd
)/../..

created_volume=()
created_network=()

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

volumes=("txchat-zookeeper" "txchat-kafka" "txchat-redis-data" "txchat-redis-config" "txchat-redis-log" "txchat-etcd-data" "txchat-comet-config" "txchat-logic-config")
networks=("txchat-components" "txchat-service")

for vname in ${volumes[*]} ; do \
  volume_create "${vname}"
done

for netname in ${networks[*]} ; do \
  network_create "${netname}"
done
