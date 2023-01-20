module github.com/txchat/im

go 1.14

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/Terry-Mao/goim v0.0.0-20210523140626-e742c99ad76e
	github.com/gin-contrib/pprof v1.3.0
	github.com/gin-gonic/gin v1.6.3
	github.com/golang/protobuf v1.5.2
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/google/uuid v1.3.0
	github.com/gorilla/websocket v1.4.2
	github.com/mitchellh/mapstructure v1.2.2
	github.com/oofpgDLD/dtask v1.0.2
	github.com/opentracing/opentracing-go v1.2.0
	github.com/rs/zerolog v1.21.0
	//github.com/txchat/dtalk v0.0.0-00010101000000-000000000000
	github.com/txchat/dtalk v0.0.1
	//github.com/txchat/im-pkg v0.0.0-00010101000000-000000000000
	github.com/txchat/im-pkg v0.0.1
	//github.com/txchat/imparse v0.0.0-00010101000000-000000000000
	github.com/txchat/imparse v0.0.1
	github.com/uber/jaeger-client-go v2.30.0+incompatible
	github.com/zeromicro/go-zero v1.4.3
	github.com/zhenjl/cityhash v0.0.0-20131128155616-cdd6a94144ab
	go.etcd.io/etcd/api/v3 v3.5.5
	go.etcd.io/etcd/client/v3 v3.5.5
	google.golang.org/genproto v0.0.0-20221111202108-142d8a6fa32e
	google.golang.org/grpc v1.50.1
	google.golang.org/protobuf v1.28.1
	gopkg.in/Shopify/sarama.v1 v1.19.0
)

replace github.com/txchat/imparse => ../imparse
replace github.com/txchat/dtalk => ../dtalk
