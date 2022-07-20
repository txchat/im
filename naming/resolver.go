package naming

import (
	"context"
	"log"
	"strings"
	"time"

	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/resolver"
)

var schema string
var cli *clientv3.Client

type etcdResolver struct {
	rawAddr string
	schema  string
	cc      resolver.ClientConn
}

func NewResolver(etcdAddr, schema string) resolver.Builder {
	return &etcdResolver{rawAddr: etcdAddr, schema: schema}
}

func (r *etcdResolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	//fmt.Println("target:", target)
	var err error
	if cli == nil {
		cli, err = clientv3.New(clientv3.Config{
			Endpoints:   strings.Split(r.rawAddr, ";"),
			DialTimeout: 15 * time.Second,
		})
		if err != nil {
			return nil, err
		}
	}

	r.cc = cc

	go r.watch("/" + target.Scheme + "/" + target.Endpoint + "/")

	return r, nil
}

func (r etcdResolver) Scheme() string {
	return r.schema
}

func (r etcdResolver) ResolveNow(rn resolver.ResolveNowOptions) {
	//log.Println("ResolveNow")
}

func (r etcdResolver) Close() {
	//log.Println("Close")
}

func (r *etcdResolver) watch(keyPrefix string) {
	var addrList []resolver.Address

	getResp, err := cli.Get(context.Background(), keyPrefix, clientv3.WithPrefix())
	if err != nil {
		log.Println(err)
	} else {
		for i := range getResp.Kvs {
			addrList = append(addrList, resolver.Address{Addr: strings.TrimPrefix(string(getResp.Kvs[i].Key), keyPrefix)})
		}
	}

	// 新版本etcd去除了NewAddress方法 以UpdateState代替
	r.cc.UpdateState(resolver.State{Addresses: addrList})

	rch := cli.Watch(context.Background(), keyPrefix, clientv3.WithPrefix())
	for n := range rch {
		for _, ev := range n.Events {
			addr := strings.TrimPrefix(string(ev.Kv.Key), keyPrefix)
			switch ev.Type {
			case mvccpb.PUT:
				if !exist(addrList, addr) {
					addrList = append(addrList, resolver.Address{Addr: addr})
					r.cc.UpdateState(resolver.State{Addresses: addrList})
				}
			case mvccpb.DELETE:
				if s, ok := remove(addrList, addr); ok {
					addrList = s
					r.cc.UpdateState(resolver.State{Addresses: addrList})
				}
			}
			log.Printf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
		}
	}
}

func exist(l []resolver.Address, addr string) bool {
	for i := range l {
		if l[i].Addr == addr {
			return true
		}
	}
	return false
}

func remove(s []resolver.Address, addr string) ([]resolver.Address, bool) {
	for i := range s {
		if s[i].Addr == addr {
			s[i] = s[len(s)-1]
			return s[:len(s)-1], true
		}
	}
	return nil, false
}
