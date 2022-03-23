package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime/pprof"
	"runtime/trace"
	"syscall"
	"time"

	gpprof "github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/oofpgDLD/dtask"
	"github.com/oofpgDLD/dtask/proto"
)

const (
	_trace     = false
	_pprof     = false
	_pprof_net = true
)

var addr = flag.String("addr", ":17070", "http service address")

func main() {
	if _trace {
		f, err := os.Create("trace.out")
		if err != nil {
			panic(err)
		}
		defer f.Close()

		err = trace.Start(f)
		if err != nil {
			panic(err)
		}
		defer trace.Stop()
	}

	if _pprof {
		f, err := os.Create("pprof.out")
		if err != nil {
			panic(err)
		}
		err = pprof.StartCPUProfile(f)
		if err != nil {
			panic(err)
		}
		defer pprof.StopCPUProfile()
	}

	r := gin.Default()
	// websocket
	r.GET("/test", ServeWs)
	r.GET("/test2", func(context *gin.Context) {
		context.String(http.StatusOK, "success")
	})

	srv := &http.Server{
		Addr:    *addr,
		Handler: r,
	}

	if _pprof_net {
		log.Printf("serve %s", *addr)
		gpprof.Register(r)
	}
	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			log.Printf("http server stop %v", err)
		}
	}()
	// init signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Printf("service get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			log.Print("server exit")
			return
		case syscall.SIGHUP:
			// TODO reload
		default:
			return
		}
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	// allow cross-origin
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

func ServeWs(context *gin.Context) {
	var (
		ck       *dtask.Clock
		jobCache map[string]dtask.Job
		wCh      chan []byte
	)
	//	userId = context.GetHeader("FZM-UID")
	c, err := upgrader.Upgrade(context.Writer, context.Request, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	wCh = make(chan []byte)
	defer c.Close()
	defer onClose(ck)
	done := make(chan struct{})
	ck = dtask.NewClock()
	jobCache = make(map[string]dtask.Job)

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
			onRecvMsg(ck, jobCache, wCh, message)
		}
	}()

	for {
		select {
		case <-done:
			return
		case data := <-wCh:
			err = c.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				log.Println("write:", err)
				break
			}
		}
	}
}

func onRecvMsg(ck *dtask.Clock, jc map[string]dtask.Job, wCh chan []byte, msg []byte) {
	msg = bytes.TrimSpace(bytes.Replace(msg, newline, space, -1))
	p := proto.Proto{}
	json.Unmarshal(msg, &p)
	switch p.Opt {
	case proto.Start:
		job, inserted := ck.AddJobRepeat(time.Second*5, 0, func() {
			wCh <- msg
		})
		if !inserted {
			break
		}
		jc[fmt.Sprint(p.Seq)] = job
	case proto.Stop:
		if job := jc[fmt.Sprint(p.Seq)]; job != nil {
			job.Cancel()
		}
	default:
	}
}

func onClose(ck *dtask.Clock) {
	if ck == nil {
		log.Printf("onClose WsConnection args can not find")
		return
	}
	ck.Stop()
}
