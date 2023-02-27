# 通讯协议说明

## 连接方式

- [X]  Websocket
- [X]  TCP

### 访问路径

```
ws://DOMAIN/sub/
tcp://DOMAIN
```

## 协议描述

> 协议包含20字节的协议头和n字节的由协议头指定长度的协议体组成。传输过程的数据为二进制的字节流。

一条完整的消息格式定义如下：


| 参数名         | 必选  | 类型            | 说明                                |
| :--------------- | :------ | :---------------- | :------------------------------------ |
| package length | true  | int32 bigendian | 包长度                              |
| header Length  | true  | int16 bigendian | 包头长度                            |
| ver            | true  | int16 bigendian | 协议版本                            |
| operation      | true  | int32 bigendian | 协议指令                            |
| seq            | true  | int32 bigendian | 序列号                              |
| ack            | true  | int32 bigendian | 确认号                              |
| body           | false | binary          | protobuf序列化字节数据 |

### operation枚举类型


| 指令 | 说明               |
| :----- | :------------------- |
| 1    | 客户端连接认证     |
| 2    | 客户端连接认证响应 |
| 3    | 客户端请求心跳     |
| 4    | 服务端心跳答复     |
| 5    | 客户端断开连接     |
| 6    | 客户端断开连接响应 |
| 7    | 发送/接收消息      |
| 8    | 接收消息响应       |

### body描述
> 不同的operation类型对应的body都需要序列化为相应的结构

#### 连接认证
operation 为 1 时， body 必须可以反序列化为 AuthMsg
```protobuf
syntax = "proto3";

message AuthMsg {
    string appId = 1;
    string token = 2;
    bytes ext = 3; // 其它业务方可能需要的信息
}
```
operation 为 2,7,8 时， body 的反序列化数据结构由具体的业务服务提供

具体协议参考 [protocol.proto](../api/protocol/protocol.proto)