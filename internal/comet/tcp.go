package comet

import (
	"bufio"
	"io"
	"net"

	"github.com/txchat/im/api/protocol"
)

type TCP struct {
	rwc io.ReadWriteCloser
	rb  *bufio.Reader
	wb  *bufio.Writer
}

func NewTCP(conn net.Conn, rb *bufio.Reader, wb *bufio.Writer) (ProtoReaderWriterCloser, error) {
	return &TCP{
		rwc: conn,
		rb:  rb,
		wb:  wb,
	}, nil
}

func (tcp *TCP) SchemeName() string {
	return "tcp"
}

func (tcp *TCP) WriteProto(p *protocol.Proto) error {
	return p.WriteTo(tcp.wb)
}

func (tcp *TCP) ReadProto(p *protocol.Proto) error {
	return p.ReadFrom(tcp.rb)
}

func (tcp *TCP) Flush() error {
	return tcp.wb.Flush()
}

func (tcp *TCP) Close() error {
	return tcp.rwc.Close()
}
