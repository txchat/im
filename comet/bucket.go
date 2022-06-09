package comet

import (
	"sync"
	"sync/atomic"

	"github.com/txchat/im/api/comet/grpc"
	"github.com/txchat/im/comet/conf"
)

// Bucket is a channel holder.
type Bucket struct {
	c     *conf.Bucket
	cLock sync.RWMutex        // protect the channels for chs
	chs   map[string]*Channel // map sub key to a channel
	//group
	groups      map[string]*Group
	routines    []chan *grpc.BroadcastGroupReq
	routinesNum uint64
}

// NewBucket new a bucket struct. store the key with im channel.
func NewBucket(c *conf.Bucket) (b *Bucket) {
	b = new(Bucket)
	b.chs = make(map[string]*Channel, c.Channel)
	b.c = c
	b.groups = make(map[string]*Group, c.Groups)
	b.routines = make([]chan *grpc.BroadcastGroupReq, c.RoutineAmount)
	for i := uint64(0); i < c.RoutineAmount; i++ {
		c := make(chan *grpc.BroadcastGroupReq, c.RoutineSize)
		b.routines[i] = c
		go b.groupProc(c)
	}
	return
}

// ChannelCount channel count in the bucket
func (b *Bucket) ChannelCount() int {
	return len(b.chs)
}

func (b *Bucket) ChannelsCount() (res []string) {
	var (
		key string
		i   int
		//ch   *Channel
	)
	b.cLock.RLock()
	res = make([]string, len(b.chs))
	for key = range b.chs {
		res[i] = key
		i++
	}
	b.cLock.RUnlock()
	return
}

// Put put a channel according with sub key.
func (b *Bucket) Put(ch *Channel) (err error) {
	b.cLock.Lock()
	// close old channel
	if dch := b.chs[ch.Key]; dch != nil {
		dch.Close()
	}
	b.chs[ch.Key] = ch
	b.cLock.Unlock()
	return
}

// Del delete the channel by sub key.
func (b *Bucket) Del(dch *Channel) {
	var (
		ok bool
		ch *Channel
	)
	b.cLock.Lock()
	if ch, ok = b.chs[dch.Key]; ok {
		if ch == dch {
			//修改内容:获取channel下所有添加的群，并逐个在群聊中删除该channel；修改人：dld;修改时间：2021年5月8日16:36:00 c4f618a9-1c37-3459-3861-a24f26bb2d85
			for id := range ch.Groups() {
				if g := b.groups[id]; g != nil {
					g.Del(ch)
				}
			}
			//结束：c4f618a9-1c37-3459-3861-a24f26bb2d85
			delete(b.chs, ch.Key)
		}
	}
	b.cLock.Unlock()
}

// Channel get a channel by sub key.
func (b *Bucket) Channel(key string) (ch *Channel) {
	b.cLock.RLock()
	ch = b.chs[key]
	b.cLock.RUnlock()
	return
}

// Broadcast push msgs to all channels in the bucket.
func (b *Bucket) Broadcast(p *grpc.Proto, op int32) {
	var ch *Channel
	b.cLock.RLock()
	for _, ch = range b.chs {
		_, _ = ch.Push(p)
	}
	b.cLock.RUnlock()
}

// group
// GroupCount room count in the bucket
func (b *Bucket) GroupCount() int {
	return len(b.groups)
}

// GroupsCount get all group id where online number > 0.
func (b *Bucket) GroupsCount() (res map[string]int32) {
	var (
		groupID string
		group   *Group
	)
	b.cLock.RLock()
	res = make(map[string]int32)
	for groupID, group = range b.groups {
		if group.Online > 0 {
			res[groupID] = group.Online
		}
	}
	b.cLock.RUnlock()
	return
}

// Put put a group according with sub key.
func (b *Bucket) PutGroup(gid string) (group *Group, err error) {
	var ok bool
	b.cLock.Lock()
	if group, ok = b.groups[gid]; !ok {
		group = NewGroup(gid)
		b.groups[gid] = group
	}
	b.cLock.Unlock()
	return
}

// Group get a group by group id.
func (b *Bucket) Group(gid string) (group *Group) {
	b.cLock.RLock()
	group = b.groups[gid]
	b.cLock.RUnlock()
	return
}

// DelGroup delete a room by group id.
func (b *Bucket) DelGroup(group *Group) {
	b.cLock.Lock()
	delete(b.groups, group.ID)
	b.cLock.Unlock()
	group.Close()
}

// BroadcastGroup broadcast a message to specified group
func (b *Bucket) BroadcastGroup(arg *grpc.BroadcastGroupReq) {
	num := atomic.AddUint64(&b.routinesNum, 1) % b.c.RoutineAmount
	b.routines[num] <- arg
}

// group proc
func (b *Bucket) groupProc(c chan *grpc.BroadcastGroupReq) {
	for {
		arg := <-c
		if group := b.Group(arg.GroupID); group != nil {
			group.Push(arg.Proto)
		}
	}
}
