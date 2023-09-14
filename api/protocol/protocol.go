package protocol

import (
	"bufio"
	"encoding/binary"
	"errors"

	"github.com/txchat/im/internal/websocket"
)

const (
	// MaxBodySize max proto body size
	MaxBodySize = uint32(1 << 14)
)

const (
	// size
	_packSize             = 4
	_headerSize           = 2
	_verSize              = 2
	_opSize               = 4
	_seqSize              = 4
	_ackSize              = 4
	_heartSize            = 4
	_rawHeaderSize uint16 = _packSize + _headerSize + _verSize + _opSize + _seqSize + _ackSize
	_maxPackSize   uint32 = MaxBodySize + uint32(_rawHeaderSize)
	// offset
	_packOffset   = 0
	_headerOffset = _packOffset + _packSize
	_verOffset    = _headerOffset + _headerSize
	_opOffset     = _verOffset + _verSize
	_seqOffset    = _opOffset + _opSize
	_ackOffset    = _seqOffset + _seqSize
	_heartOffset  = _ackOffset + _ackSize
)

var (
	// ErrProtoPackLen proto packet len error
	ErrProtoPackLen = errors.New("default server codec pack length error")
	// ErrProtoHeaderLen proto header len error
	ErrProtoHeaderLen = errors.New("default server codec header length error")
)

var (
	// ProtoReady proto ready
	ProtoReady = &Proto{Op: int32(Op_ProtoReady)}
	// ProtoFinish proto finish
	ProtoFinish = &Proto{Op: int32(Op_ProtoFinish)}
	// ProtoResend proto resend
	ProtoResend = &Proto{Op: int32(Op_ProtoResend)}
)

func peek(w *bufio.Writer, n int) ([]byte, error) {
	var err error
	if n < 0 {
		return nil, bufio.ErrNegativeCount
	}
	if n > w.Size() {
		return nil, bufio.ErrBufferFull
	}
	for w.Available() < n && err == nil {
		err = w.Flush()
	}
	if err != nil {
		return nil, err
	}
	buf := w.AvailableBuffer()
	return buf[:n], nil
}

func pop(r *bufio.Reader, n int) ([]byte, error) {
	buf, err := r.Peek(n)
	if err != nil {
		return nil, err
	}
	if _, err = r.Discard(n); err != nil {
		return nil, err
	}
	return buf, nil
}

// WriteTo write a proto to bytes writer.
func (p *Proto) WriteTo(w *bufio.Writer) (err error) {
	var (
		buf     []byte
		packLen uint32
	)
	packLen = uint32(_rawHeaderSize) + uint32(len(p.Body))
	if buf, err = peek(w, int(_rawHeaderSize)); err != nil {
		return
	}
	binary.BigEndian.PutUint32(buf[_packOffset:], packLen)
	binary.BigEndian.PutUint16(buf[_headerOffset:], _rawHeaderSize)
	binary.BigEndian.PutUint16(buf[_verOffset:], uint16(p.Ver))
	binary.BigEndian.PutUint32(buf[_opOffset:], uint32(p.Op))
	binary.BigEndian.PutUint32(buf[_seqOffset:], uint32(p.Seq))
	binary.BigEndian.PutUint32(buf[_ackOffset:], uint32(p.Ack))
	_, err = w.Write(buf)
	if err != nil {
		return
	}
	if p.Body != nil {
		_, err = w.Write(p.Body)
	}
	return
}

func (p *Proto) ReadFrom(r *bufio.Reader) (err error) {
	var (
		packLen   uint32
		headerLen uint16
		bodyLen   int
		buf       []byte
	)
	if buf, err = pop(r, int(_rawHeaderSize)); err != nil {
		return
	}

	packLen = binary.BigEndian.Uint32(buf[_packOffset:_headerOffset])
	headerLen = binary.BigEndian.Uint16(buf[_headerOffset:_verOffset])
	p.Ver = int32(binary.BigEndian.Uint16(buf[_verOffset:_opOffset]))
	p.Op = int32(binary.BigEndian.Uint32(buf[_opOffset:_seqOffset]))
	p.Seq = int32(binary.BigEndian.Uint32(buf[_seqOffset:_ackOffset]))
	p.Ack = int32(binary.BigEndian.Uint32(buf[_ackOffset:]))
	if packLen > _maxPackSize {
		return ErrProtoPackLen
	}
	if headerLen != _rawHeaderSize {
		return ErrProtoHeaderLen
	}
	if bodyLen = int(packLen - uint32(headerLen)); bodyLen > 0 {
		p.Body, err = pop(r, bodyLen)
	} else {
		p.Body = nil
	}
	return
}

type WebsocketReaderWriter interface {
	ReadMessage() (int, []byte, error)
}

// ReadWebsocket read a proto from websocket connection.
func (p *Proto) ReadWebsocket(ws *websocket.Conn) (err error) {
	var (
		packLen   uint32
		headerLen uint16
		bodyLen   int
		buf       []byte
	)
	if _, buf, err = ws.ReadMessage(); err != nil {
		return
	}
	if len(buf) < int(_rawHeaderSize) {
		return ErrProtoPackLen
	}
	packLen = binary.BigEndian.Uint32(buf[_packOffset:_headerOffset])
	headerLen = binary.BigEndian.Uint16(buf[_headerOffset:_verOffset])
	p.Ver = int32(binary.BigEndian.Uint16(buf[_verOffset:_opOffset]))
	p.Op = int32(binary.BigEndian.Uint32(buf[_opOffset:_seqOffset]))
	p.Seq = int32(binary.BigEndian.Uint32(buf[_seqOffset:_ackOffset]))
	p.Ack = int32(binary.BigEndian.Uint32(buf[_ackOffset:]))
	if packLen > _maxPackSize {
		return ErrProtoPackLen
	}
	if headerLen != _rawHeaderSize {
		return ErrProtoHeaderLen
	}
	if bodyLen = int(packLen - uint32(headerLen)); bodyLen > 0 {
		p.Body = buf[headerLen:packLen]
	} else {
		p.Body = nil
	}
	return
}

// WriteWebsocket write a proto to websocket connection.
func (p *Proto) WriteWebsocket(ws *websocket.Conn) (err error) {
	var (
		packLen uint32
		buf     []byte
	)
	packLen = uint32(_rawHeaderSize) + uint32(len(p.Body))
	if err = ws.WriteHeader(websocket.BinaryMessage, int(packLen)); err != nil {
		return
	}
	if buf, err = ws.Peek(int(_rawHeaderSize)); err != nil {
		return
	}
	binary.BigEndian.PutUint32(buf[_packOffset:], packLen)
	binary.BigEndian.PutUint16(buf[_headerOffset:], _rawHeaderSize)
	binary.BigEndian.PutUint16(buf[_verOffset:], uint16(p.Ver))
	binary.BigEndian.PutUint32(buf[_opOffset:], uint32(p.Op))
	binary.BigEndian.PutUint32(buf[_seqOffset:], uint32(p.Seq))
	binary.BigEndian.PutUint32(buf[_ackOffset:], uint32(p.Ack))
	if p.Body != nil {
		err = ws.WriteBody(p.Body)
	}
	return
}
