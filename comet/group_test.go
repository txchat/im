package comet

import (
	"strconv"
	"sync/atomic"
	"testing"
	"time"
)

func TestGroup_Put(t *testing.T) {
	//one channel
	var gid int32 = 0
	ch := NewChannel(10, 10)
	//one groups
	{
		g := NewGroup(strconv.Itoa(int(atomic.AddInt32(&gid, 1))))
		for i := 0; i < 10; i++ {
			go func() {
				for {
					err := g.Put(ch)
					if err != nil {
						t.Error(err)
					}
				}
			}()
		}
	}

	//many groups
	{
		for i := 0; i < 10; i++ {
			g := NewGroup(strconv.Itoa(int(atomic.AddInt32(&gid, 1))))
			go func() {
				for {
					err := g.Put(ch)
					if err != nil {
						t.Error(err)
					}
				}
			}()
		}
	}

	time.Sleep(10 * time.Second)
}

func TestGroup_Del(t *testing.T) {
	//one channel
	var gid int32 = 0
	//var lc sync.RWMutex
	var groups = make(map[int32]*Group)
	ch := NewChannel(10, 10)
	//init
	for i := 0; i < 10; i++ {
		g := NewGroup(strconv.Itoa(int(atomic.AddInt32(&gid, 1))))
		groups[gid] = g
	}
	gid = 0

	go func() {
		for _, g := range groups {
			err := g.Put(ch)
			if err != nil {
				t.Error(err)
			}
		}
	}()

	go func() {
		for _, g := range groups {
			g.Del(ch)
		}
	}()
	time.Sleep(10 * time.Second)
}

func TestGroup_Del2(t *testing.T) {
	//one channel
	var gid int32 = 0
	//var lc sync.RWMutex
	var groups = make(map[int32]*Group)
	ch := NewChannel(10, 10)
	//init
	for i := 0; i < 10; i++ {
		g := NewGroup(strconv.Itoa(int(atomic.AddInt32(&gid, 1))))
		groups[gid] = g
	}
	gid = 0

	go func() {
		for _, g := range groups {
			err := g.Put(ch)
			if err != nil {
				t.Error(err)
			}
		}
	}()

	go func() {
		for _, g := range groups {
			g.Del(ch)
		}
	}()
	time.Sleep(10 * time.Second)
}

func convertAddress(addr string) string {
	if len(addr) == 0 {
		return addr
	}
	switch addr[0] {
	case 'g':
		return addr[7:]
	default:
		return addr
	}
}

func Test_convertAddress(t *testing.T) {
	t.Log(convertAddress("grpc://172.16.101.107:3109"))
}
