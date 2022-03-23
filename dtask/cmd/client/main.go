package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
	"github.com/oofpgDLD/dtask/proto"
)

var addr = flag.String("addr", "172.16.101.107:17070", "http service address")

var (
	closer = make(chan int, 1)
	times  = 5
	users  = 10000

	u = url.URL{Scheme: "ws", Host: *addr, Path: "/test"}
)

func main() {
	flag.Parse()
	log.SetFlags(0)

	for i := 0; i < users; i++ {
		go NewUser(i, times)
	}

	go getInputByScanner()
	// init signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Printf("client get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			log.Print("client exit")
			return
		case syscall.SIGHUP:
			// TODO reload
		default:
			return
		}
	}
}

func NewUser(id int, times int) {
	uid := fmt.Sprint(id)
	log.Printf("connecting to %s, id %v", u.String(), fmt.Sprint(id))
	c, resp, err := websocket.DefaultDialer.Dial(u.String(), http.Header{
		"FZM-UID": {fmt.Sprint(id)},
	})
	if err != nil {
		log.Fatal("dial:", err, resp, u.String())
		return
	}
	defer c.Close()
	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()
	seq := 0
	ticker := time.NewTicker(time.Millisecond * 500)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			seq++
			msg, err := Msg(uid, int64(seq))
			if err != nil {
				log.Println("marshal proto:", err)
				return
			}
			err = c.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Println("write:", err)
				return
			}
			if seq >= times {
				ticker.Stop()
			}
		case index := <-closer:
			log.Println("interrupt")

			if index == id {
				log.Println("exit")
				return
			}
			closer <- index
		}
	}
}

func Msg(uid string, seq int64) ([]byte, error) {
	p := &proto.Proto{
		Uid: uid,
		Opt: proto.Start,
		Seq: seq,
	}
	return json.Marshal(p)
}

func getInputByScanner() string {
	var str string
	for {
		//使用os.Stdin开启输入流
		in := bufio.NewScanner(os.Stdin)
		if in.Scan() {
			str = in.Text()
		} else {
			str = "Find input error"
		}

		index, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			continue
		}
		closer <- int(index)
	}
}
