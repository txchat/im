package dao

import (
	"os"
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/txchat/im/app/logic/internal/config"
)

var (
	testConf  *config.Redis
	testRedis *redis.Pool
)

func TestMain(m *testing.M) {
	testRedis = newRedis(&config.Redis{
		Network:      "tcp",
		Addr:         "127.0.0.1:6379",
		Auth:         "",
		Active:       60000,
		Idle:         1024,
		DialTimeout:  200 * time.Millisecond,
		ReadTimeout:  500 * time.Millisecond,
		WriteTimeout: 500 * time.Millisecond,
		IdleTimeout:  120 * time.Second,
		Expire:       30 * time.Minute,
	})
	os.Exit(m.Run())
}
