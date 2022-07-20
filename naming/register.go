package naming

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func Register(etcdAddr, name, addr, schema string, ttl int64) error {
	var err error
	if cli == nil {
		cli, err = clientv3.New(clientv3.Config{
			Endpoints:   strings.Split(etcdAddr, ";"),
			DialTimeout: 5 * time.Second,
		})
		if err != nil {
			log.Printf("connect to etcd err:%s", err)
			return err
		}
	}

	ticker := time.NewTicker(time.Second * time.Duration(ttl))

	go func() {
		for {
			getResp, err := cli.Get(context.Background(), "/"+schema+"/"+name+"/"+addr)
			if err != nil {
				fmt.Println("etcd get err:", err)
			} else if getResp.Count == 0 {
				err = withAlive(name, addr, schema, ttl)
				if err != nil {
					log.Printf("keep alive:%s", err)
				}
			}
			<-ticker.C
		}
	}()

	return nil
}

func withAlive(name, addr, schema string, ttl int64) error {
	leaseResp, err := cli.Grant(context.Background(), ttl)
	if err != nil {
		return err
	}

	log.Printf("key:%v\n", "/"+schema+"/"+name+"/"+addr)
	_, err = cli.Put(context.Background(), "/"+schema+"/"+name+"/"+addr, addr, clientv3.WithLease(leaseResp.ID))
	if err != nil {
		log.Printf("put etcd error:%s", err)
		return err
	}

	ch, err := cli.KeepAlive(context.Background(), leaseResp.ID)
	if err != nil {
		log.Printf("keep alive error:%s", err)
		return err
	}

	go func() {
		for range ch {
		}
		//for leaseKeepResp := range ch {
		//	log.Println("续约成功", leaseKeepResp.ID)
		//}
		//
		//log.Println("关闭续租")
	}()

	return nil
}

func UnRegister(name, addr, schema string) error {
	if cli != nil {
		//fmt.Println("unregister...")
		_, err := cli.Delete(context.Background(), "/"+schema+"/"+name+"/"+addr)
		return err
	}

	return nil
}
