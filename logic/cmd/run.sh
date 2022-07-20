#!/bin/bash

# -log_dir=log 日志输出文件夹
# -logtostderr=true 打印到标准错误而不是文件。
# -alsologtostderr=true 同时打印到标准错误。
# -v=4 设置日志级别
# -vmodule=main=5 设置单个文件日志级别
./main -v=4 -log_dir=log -alsologtostderr=true
