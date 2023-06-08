package config

import (
	xkafka "github.com/oofpgDLD/kafka-go"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	AppID             string `json:",default=echo"`
	LogicRPC          zrpc.RpcClientConf
	Producer          xkafka.ProducerConfig
	ConsumerConfig    xkafka.ConsumerConfig
	BatchConsumerConf xkafka.BatchConsumerConf
}
