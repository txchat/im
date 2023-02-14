package main

// Start Commond eg: ./client 1 1000 localhost:3102
// first parameter：beginning userId
// second parameter: amount of clients
// third parameter: comet ws-server addr

import (
	"bufio"
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/txchat/im/api/protocol"
)

const (
	rawHeaderLen = uint16(16)
	heart        = 5 * time.Second
)

type Proto struct {
	PackLen   int32  // package length
	HeaderLen int16  // header length
	Ver       int16  // protocol version
	Op        int32  // operation for request
	Seq       int32  // sequence number chosen by client
	Body      []byte // body
}

var (
	countDown  int64
	aliveCount int64
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()
	begin, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}
	num, err := strconv.Atoi(os.Args[2])
	if err != nil {
		panic(err)
	}
	go result()
	for i := begin; i < begin+num; i++ {
		go client(int64(i))
	}
	// signal
	var exit chan bool
	<-exit
}

func result() {
	var (
		lastTimes int64
		interval  = int64(5)
	)
	for {
		nowCount := atomic.LoadInt64(&countDown)
		nowAlive := atomic.LoadInt64(&aliveCount)
		diff := nowCount - lastTimes
		lastTimes = nowCount
		fmt.Println(fmt.Sprintf("%s alive:%d down:%d down/s:%d", time.Now().Format("2006-01-02 15:04:05"), nowAlive, nowCount, diff/interval))
		time.Sleep(time.Second * time.Duration(interval))
	}
}

func client(uid int64) {
	for {
		startClient(uid)
		time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
	}
}

func startClient(key int64) {
	time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
	atomic.AddInt64(&aliveCount, 1)
	quit := make(chan bool, 1)
	defer func() {
		close(quit)
		atomic.AddInt64(&aliveCount, -1)
	}()

	// connnect to server
	wsUrl := "ws://" + os.Args[3] + "/sub"
	fmt.Println("wsUrl", wsUrl)
	conn, _, err := websocket.DefaultDialer.Dial(wsUrl, nil)
	if err != nil {
		panic(err)
	}

	seq := int32(0)
	authMsg := &protocol.AuthBody{
		AppId: "echo",
		Token: "fdasfdsaf",
	}
	p := new(Proto)
	p.Ver = 1
	p.Op = int32(protocol.Op_Auth)
	p.Seq = seq
	p.Body, _ = proto.Marshal(authMsg)

	//auth
	if err = wsWriteProto(conn, p); err != nil {
		log.Printf("wsWriteProto() error(%v)", err)
		return
	}
	if err = wsReadProto(conn, p); err != nil {
		log.Printf("tcpReadProto() error(%v)", err)
		return
	}
	log.Printf("key:%d auth ok, proto: %v", key, p)

	seq++
	// writer
	go func() {
		hbProto := new(Proto)
		for {
			// heartbeat
			hbProto.Op = int32(protocol.Op_Heartbeat)
			hbProto.Seq = seq
			hbProto.Body = nil
			if err = wsWriteProto(conn, hbProto); err != nil {
				log.Printf("key:%d tcpWriteProto() error(%v)", key, err)
				return
			}
			log.Printf("key:%d Write heartbeat", key)
			time.Sleep(heart)
			seq++
			select {
			case <-quit:
				return
			default:
			}
		}
	}()
	// reader
	for {
		if err = wsReadProto(conn, p); err != nil {
			log.Printf("key:%d tcpReadProto() error(%v)", key, err)
			quit <- true
			return
		}
		if p.Op == int32(protocol.Op_AuthReply) {
			log.Printf("key:%d auth success", key)
		} else if p.Op == int32(protocol.Op_HeartbeatReply) {
			log.Printf("key:%d receive heartbeat reply", key)
			if err = conn.SetReadDeadline(time.Now().Add(heart + 60*time.Second)); err != nil {
				log.Printf("conn.SetReadDeadline() error(%v)", err)
				quit <- true
				return
			}
		} else {
			log.Printf("key:%d op:%d msg: %s", key, p.Op, string(p.Body))
			atomic.AddInt64(&countDown, 1)
		}
	}
}

func wsWriteProto(conn *websocket.Conn, proto *Proto) (err error) {
	wc, err := conn.NextWriter(websocket.BinaryMessage)
	if err != nil {
		panic(err)
	}
	wr := bufio.NewWriter(wc)

	// write
	if err = binary.Write(wr, binary.BigEndian, uint32(rawHeaderLen)+uint32(len(proto.Body))); err != nil {
		return
	}
	if err = binary.Write(wr, binary.BigEndian, rawHeaderLen); err != nil {
		return
	}
	if err = binary.Write(wr, binary.BigEndian, proto.Ver); err != nil {
		return
	}
	if err = binary.Write(wr, binary.BigEndian, proto.Op); err != nil {
		return
	}
	if err = binary.Write(wr, binary.BigEndian, proto.Seq); err != nil {
		return
	}
	if proto.Body != nil {
		if err = binary.Write(wr, binary.BigEndian, proto.Body); err != nil {
			return
		}
	}
	err = wr.Flush()
	wc.Close()
	return
}

func wsReadProto(conn *websocket.Conn, proto *Proto) (err error) {
	var (
		packLen   int32
		headerLen int16
	)

	_, rc, err := conn.NextReader()
	if err != nil {
		log.Printf("NextReader error(%v)", err)
		return err
	}
	rd := bufio.NewReader(rc)

	// read
	if err = binary.Read(rd, binary.BigEndian, &packLen); err != nil {
		return
	}
	if err = binary.Read(rd, binary.BigEndian, &headerLen); err != nil {
		return
	}
	if err = binary.Read(rd, binary.BigEndian, &proto.Ver); err != nil {
		return
	}
	if err = binary.Read(rd, binary.BigEndian, &proto.Op); err != nil {
		return
	}
	if err = binary.Read(rd, binary.BigEndian, &proto.Seq); err != nil {
		return
	}
	var (
		n, t    int
		bodyLen = int(packLen - int32(headerLen))
	)
	if bodyLen > 0 {
		proto.Body = make([]byte, bodyLen)
		for {
			if t, err = rd.Read(proto.Body[n:]); err != nil {
				return
			}
			if n += t; n == bodyLen {
				break
			}
		}
	} else {
		proto.Body = nil
	}
	return
}
