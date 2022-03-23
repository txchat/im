package dao

import (
	"os"
	"testing"
	"time"

	xtime "github.com/Terry-Mao/goim/pkg/time"
	"github.com/gomodule/redigo/redis"
	"github.com/txchat/im/logic/conf"
)

var (
	testConf  *conf.Config
	testRedis *redis.Pool
)

func TestMain(m *testing.M) {
	testRedis = newRedis(&conf.Redis{
		Network:      "tcp",
		Addr:         "127.0.0.1:6379",
		Auth:         "",
		Active:       60000,
		Idle:         1024,
		DialTimeout:  xtime.Duration(200 * time.Millisecond),
		ReadTimeout:  xtime.Duration(500 * time.Millisecond),
		WriteTimeout: xtime.Duration(500 * time.Millisecond),
		IdleTimeout:  xtime.Duration(120 * time.Second),
		Expire:       xtime.Duration(30 * time.Minute),
	})
	os.Exit(m.Run())
}
