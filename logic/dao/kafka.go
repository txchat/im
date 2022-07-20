package dao

import (
	"context"
	"fmt"

	"github.com/opentracing/opentracing-go"
	"github.com/txchat/im-pkg/trace"
	comet "github.com/txchat/im/api/comet/grpc"

	"github.com/golang/protobuf/proto"
	"github.com/rs/zerolog/log"
	pb "github.com/txchat/im/api/logic/grpc"
	"gopkg.in/Shopify/sarama.v1"
)

// PushMsg push a message to databus.
func (d *Dao) PublishMsg(ctx context.Context, appId string, fromId string, op comet.Op, key string, msg []byte) (err error) {
	tracer := opentracing.GlobalTracer()
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, tracer, fmt.Sprintf("Publish -%v-%v", appId, op.String()))
	defer span.Finish()

	pushMsg := &pb.BizMsg{
		AppId:  appId,
		FromId: fromId,
		Op:     int32(op),
		Key:    key,
		Msg:    msg,
	}
	b, err := proto.Marshal(pushMsg)
	if err != nil {
		return
	}
	appTopic := fmt.Sprintf("goim-%s-topic", appId)
	m := &sarama.ProducerMessage{
		Key:   sarama.StringEncoder(fromId),
		Topic: appTopic,
		Value: sarama.ByteEncoder(b),
	}
	trace.InjectMQHeader(tracer, span.Context(), ctx, m)
	if _, _, err = d.kafkaPub.SendMessage(m); err != nil {
		log.Error().Interface("pushMsg", pushMsg).Err(err).Msg("kafkaPub.SendMessage error")
	}
	return
}
