package comet

import (
	"context"
	"strconv"

	"github.com/txchat/im/api/comet/grpc"
	"github.com/txchat/im/dtask"
)

// Operate operate.
func (s *Comet) Operate(ctx context.Context, p *grpc.Proto, ch *Channel, tsk *dtask.Task) error {
	switch p.Op {
	case int32(grpc.Op_SendMsg):
		//标明Ack的消息序列
		p.Ack = p.Seq
		err := s.Receive(ctx, ch.Key, p)
		if err != nil {
			//下层业务调用失败，返回error的话会直接断开连接
			return err
		}
		//p.Op = int32(grpc.Op_SendMsgReply)
	case int32(grpc.Op_ReceiveMsgReply):
		//从task中删除某一条
		if j := tsk.Get(strconv.FormatInt(int64(p.Ack), 10)); j != nil {
			j.Cancel()
		}
		err := s.Receive(ctx, ch.Key, p)
		if err != nil {
			//下层业务调用失败，返回error的话会直接断开连接
			return err
		}
	case int32(grpc.Op_SyncMsgReq):
		err := s.Receive(ctx, ch.Key, p)
		if err != nil {
			//下层业务调用失败，返回error的话会直接断开连接
			return err
		}
		p.Op = int32(grpc.Op_SyncMsgReply)
	default:
		return s.Receive(ctx, ch.Key, p)
	}
	return nil
}
