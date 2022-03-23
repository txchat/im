# comet 客户端通讯协议    
                                               
`comet` 支持两种协议和客户端通讯 `websocket`，`tcp`。
                                                             
**请求URL**

```
ws://DOMAIN/sub/
tcp://DOMAIN
```

**协议格式**

二进制，请求和返回协议一致

**请求&返回参数**

| 参数名     | 必选  | 类型 | 说明       |
| :-----     | :---  | :--- | :---       |
| package length        | true  | int32 bigendian | 包长度 |
| header Length         | true  | int16 bigendian    | 包头长度 |
| ver        | true  | int16 bigendian    | 协议版本 |
| operation          | true | int32 bigendian | 协议指令 |
| seq         | true | int32 bigendian | 序列号 |
| ack         | true | int32 bigendian | 确认号 |
| body         | false | binary | $(package lenth) - $(header length) |

## 指令
| 指令     | 说明  | 
| :-----     | :---  |
| 1 | 客户端连接认证 |
| 2 | 客户端连接认证响应 |
| 3 | 客户端请求心跳 |
| 4 | 服务端心跳答复 |
| 5 | 客户端断开连接 |
| 6 | 客户端断开连接响应 |
| 7 | 客户端发送消息 |
| 8 | 客户端发送消息响应 |
| 9 | 客户端接收消息 |
| 10 | 客户端接收消息响应 |
| 14 | 多端同步指令（未使用） |
| 15 | 多端同步指令响应（未使用） |


### 连接认证

operation 为 1 时， body 必须可以反序列化为 AuthMsg

```
message AuthMsg {
    string appId = 1;
    string token = 2;
    bytes ext = 3; // 其它业务方可能需要的信息
}
```

operation 为 15 时， body 必须可以反序列化为 SyncMsg
```
message SyncMsg {
    int64 logId = 1;
}
```
具体协议参考 `api/comet/grpc/comet.proto` 和 `api/logic/grpc/comet.proto`