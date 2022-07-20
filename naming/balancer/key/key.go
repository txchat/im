package key

import (
	"sync"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/grpclog"
)

const (
	Name       = "picker_key"
	DefaultKey = "X-rand-string-balabala"
)

func init() {
	balancer.Register(newBuilder(DefaultKey))
}

func newBuilder(key string) balancer.Builder {
	return base.NewBalancerBuilder(
		Name,
		&keyPickerBuilder{key: key},
		base.Config{HealthCheck: false},
	)
}

type keyPickerBuilder struct {
	key string
}

func (b *keyPickerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	grpclog.Infof("keyPickerBuilder: newPicker called with info: %v", info)
	if len(info.ReadySCs) == 0 {
		return base.NewErrPicker(balancer.ErrNoSubConnAvailable)
	}

	picker := &keyPicker{
		addr2sc: make(map[string]balancer.SubConn),
		key:     b.key,
	}

	for sc, conInfo := range info.ReadySCs {
		addr := conInfo.Address.Addr
		picker.setAddr2sc(addr, sc)
	}

	return picker
}

type keyPicker struct {
	addr2sc map[string]balancer.SubConn
	sync.RWMutex
	key string
}

func (p *keyPicker) setAddr2sc(addr string, sc balancer.SubConn) {
	p.Lock()
	defer p.Unlock()
	p.addr2sc[addr] = sc
	log.Debug().Str("addr", addr).Msg("set addr 2 sc")
}

func (p *keyPicker) getAddr2sc(addr string) (balancer.SubConn, error) {
	p.RLock()
	defer p.RUnlock()
	sc, ok := p.addr2sc[addr]
	if ok {
		return sc, nil
	}
	log.Warn().Str("addr", addr).Msg("server warn")
	// TODO 如果找不到就随便给一个
	for _, tsc := range p.addr2sc {

		return tsc, nil
	}
	return nil, nil
}

func (p *keyPicker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	var res balancer.PickResult

	addr, _ := info.Ctx.Value(p.key).(string)

	sc, _ := p.getAddr2sc(addr)
	res.SubConn = sc
	return res, nil
}
