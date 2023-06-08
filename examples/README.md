# 网关使用示例

本示例实现一个最简单的**回声服务**（将客户端发送的Ping消息按原路返回一个Pong消息）来展示如何如何使用消息网关。

## 实现业务服务

### 实现鉴权插件

消息网关会根据握手消息中的AppID访问相应的鉴权插件，`echo`服务需要实现该插件接口[pkg/auth/auth.go](../pkg/auth/auth.go)：

`echo`在实现鉴权插件是并没有通过网络访问业务服务的接口，只是简单的从token中提取UID并返回，在实际开发中可以通过HTTP访问。
- [server/auth.go](server/auth.go)
- [server/auth.go](server/auth.go)

logic服务的`main`包导入`import _ "github.com/txchat/im/examples/server"`，并在服务配置文件添加插件列表：
```yaml
Apps:
  -
    AppId: "mock"
    AuthURL: ""
    Timeout: "5s"
```

### 实现处理逻辑

客户端进行连接的建立断开，消息推送和ACK等事件最终会在消息中间件中产生数据流，业务服务可以监听消息中间件进行消费，不同的事件对应协议中不同的操作类型，具体参看[协议格式protocol](../api/protocol/protocol.proto)

`echo`服务在[echo.proto](server/echo/types/echo.proto)文件定义了Ping和Pong的协议格式，并在[mq.go](server/echo/mq/mq.go)中处理网关协议的`Op=Message`类型的消息，解开消息体并按原路返回对应Pong消息。

## 编译打包

### 从源码编译

环境要求: 
1. Golang 1.17 or later, 参考[golang官方安装文档](https://go.dev/doc/install)

```shell
# 编译本机系统和指令集的可执行文件
$ make build

# 编译目标机器的可执行文件,例如
$ make build_linux_amd64
```
编译成功后目标执行文件在工程目录的target文件夹下。

运行服务命令：
```shell
$ ./target/comet -f ./target/comet.toml
$ ./target/logic -f ./target/logic.toml
```

### docker容器中运行

环境要求:
1. Golang 1.17 or later, 参考[golang官方安装文档](https://go.dev/doc/install)
2. docker engine version 20.10.17 or later, 安装参考[docker官方安装文档](https://docs.docker.com/get-docker/)

```shell
# 初始化docker环境
$ make init-compose

# 打包镜像及运行容器
$ make docker-compose-up

# 查看容器是否运行成功
$ make docker-compose-ps
```

## 测试

[im-util](https://github.com/txchat/im-util)仓库下有一系列工具，找到该仓库下`app-examples/echo`目录下的客户端测试工具编译运行。