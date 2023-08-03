package comet

import (
	"io"

	"github.com/Terry-Mao/goim/pkg/bufio"
	"github.com/txchat/im/api/protocol"
)

type TCP struct {
	rwc io.ReadWriteCloser
	rr  *bufio.Reader
	wr  *bufio.Writer
}

func NewTCP(rwc io.ReadWriteCloser) *TCP {
	return &TCP{
		rwc: rwc,
	}
}

func (tcp *TCP) SchemeName() string {
	return "tcp"
}

func (tcp *TCP) Upgrade(rr *bufio.Reader, wr *bufio.Writer) error {
	tcp.rr = rr
	tcp.wr = wr
	return nil
}

func (tcp *TCP) WriteProto(p *protocol.Proto) error {
	return p.WriteTCP(tcp.wr)
}

func (tcp *TCP) WriteHeart(p *protocol.Proto, online int32) error {
	return p.WriteTCPHeart(tcp.wr, online)
}

func (tcp *TCP) ReadProto(p *protocol.Proto) error {
	return p.ReadTCP(tcp.rr)
}

func (tcp *TCP) Flush() error {
	return tcp.wr.Flush()
}

func (tcp *TCP) Close() error {
	return tcp.rwc.Close()
}
