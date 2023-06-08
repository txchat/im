package mq

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/proto"
	xkafka "github.com/oofpgDLD/kafka-go"
	"github.com/oofpgDLD/kafka-go/trace"
	"github.com/txchat/im/api/protocol"
	"github.com/txchat/im/app/logic/logicclient"
	"github.com/txchat/im/examples/server/echo/internal/config"
	"github.com/txchat/im/examples/server/echo/internal/model"
	"github.com/txchat/im/examples/server/echo/internal/svc"
	"github.com/txchat/im/examples/server/echo/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type Service struct {
	logx.Logger
	Config        config.Config
	svcCtx        *svc.ServiceContext
	batchConsumer *xkafka.BatchConsumer
}

func NewService(cfg config.Config, svcCtx *svc.ServiceContext) *Service {
	s := &Service{
		Logger: logx.WithContext(context.TODO()),
		Config: cfg,
		svcCtx: svcCtx,
	}
	//topic config
	cfg.ConsumerConfig.Topic = fmt.Sprintf("goim-%s-receive", cfg.AppID)
	cfg.ConsumerConfig.Group = fmt.Sprintf("goim-%s-receive-transfer", cfg.AppID)
	//new batch consumer
	consumer := xkafka.NewConsumer(cfg.ConsumerConfig, nil)
	logx.Info("dial kafka broker success")
	bc := xkafka.NewBatchConsumer(cfg.BatchConsumerConf, consumer, xkafka.WithHandle(s.handleFunc), xkafka.WithBatchConsumerInterceptors(trace.ConsumeInterceptor))
	s.batchConsumer = bc
	return s
}

func (s *Service) Serve() {
	s.batchConsumer.Start()
}

func (s *Service) Shutdown(ctx context.Context) {
	s.batchConsumer.GracefulStop(ctx)
}

func (s *Service) handleFunc(ctx context.Context, key string, data []byte) error {
	var receivedMsg logicclient.ReceivedMessage
	if err := proto.Unmarshal(data, &receivedMsg); err != nil {
		s.Error("proto.Unmarshal receivedMessage error", "err", err)
		return err
	}
	if receivedMsg.GetAppId() != s.Config.AppID {
		s.Error(model.ErrAppID.Error())
		return model.ErrAppID
	}

	switch receivedMsg.GetOp() {
	case protocol.Op_Message:
		var echo types.EchoMsg
		err := proto.Unmarshal(receivedMsg.GetBody(), &echo)
		if err != nil {
			return err
		}

		receivedMsg.Body, err = makeResp(&echo)
		if err != nil {
			return err
		}

		bytes, err := proto.Marshal(&receivedMsg)
		if err != nil {
			return err
		}

		_, err = s.svcCtx.LogicClient.PushByKey(ctx, &logicclient.PushByKeyReq{
			AppId: s.svcCtx.Config.AppID,
			ToKey: []string{receivedMsg.Key},
			Op:    receivedMsg.GetOp(),
			Body:  bytes,
		})
		if err != nil {
			return err
		}
	default:
		return model.ErrCustomNotSupport
	}
	return nil
}

func makeResp(e *types.EchoMsg) ([]byte, error) {
	var ee types.EchoMsg
	ee.Ty = int32(types.EchoOp_PongAction)
	ee.Value = &types.EchoMsg_Pong{Pong: &types.Pong{Msg: fmt.Sprintf("pong for %s", e.Value.(*types.EchoMsg_Ping).Ping.Msg)}}
	return proto.Marshal(&ee)
}
