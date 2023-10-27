package comet

import (
	"sync"
	"sync/atomic"

	"github.com/txchat/im/api/protocol"
)

// BucketConfig is bucket config.
type BucketConfig struct {
	Size          int    `json:",default=32"`
	Channel       int    `json:",default=1024"`
	Groups        int    `json:",default=1024"`
	RoutineAmount uint64 `json:",default=32"`
	RoutineSize   int    `json:",default=1024"`
}

type GroupCastReq struct {
	gid   string
	proto *protocol.Proto
}

// Bucket is a channel holder.
type Bucket struct {
	cfg   *BucketConfig
	cLock sync.RWMutex        // protect the channels for chs
	chs   map[string]*Channel // map sub key to a channel
	//group
	groups      map[string]*Group
	routines    []chan *GroupCastReq
	routinesNum uint64
}

// NewBucket new a bucket struct. store the key with im channel.
func NewBucket(cfg *BucketConfig) (b *Bucket) {
	b = new(Bucket)
	b.cfg = cfg
	b.chs = make(map[string]*Channel, cfg.Channel)
	b.groups = make(map[string]*Group, cfg.Groups)
	b.routines = make([]chan *GroupCastReq, cfg.RoutineAmount)
	for i := uint64(0); i < cfg.RoutineAmount; i++ {
		groupChan := make(chan *GroupCastReq, cfg.RoutineSize)
		b.routines[i] = groupChan
		go b.groupProc(groupChan)
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

// Put hold a channel instance.
func (b *Bucket) Put(ch *Channel) {
	b.cLock.Lock()
	// close old channel
	if dch := b.chs[ch.Key]; dch != nil {
		dch.Close()
	}
	b.chs[ch.Key] = ch
	b.cLock.Unlock()
}

// Del delete the channel instance.
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

// Broadcast push messages to all channels in the bucket.
func (b *Bucket) Broadcast(p *protocol.Proto) {
	var ch *Channel
	b.cLock.RLock()
	for _, ch = range b.chs {
		_ = ch.Push(p)
	}
	b.cLock.RUnlock()
}

// GroupCount groups count in the bucket
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

// PutGroup new and got group instance by id.
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

// DelGroup delete group instance.
func (b *Bucket) DelGroup(group *Group) {
	b.cLock.Lock()
	delete(b.groups, group.ID)
	b.cLock.Unlock()
	group.Close()
}

// GroupCast broadcast a message in the specified group
func (b *Bucket) GroupCast(gid string, p *protocol.Proto) {
	num := atomic.AddUint64(&b.routinesNum, 1) % b.cfg.RoutineAmount
	b.routines[num] <- &GroupCastReq{
		gid:   gid,
		proto: p,
	}
}

// groupProc goroutine for process group push request
func (b *Bucket) groupProc(c chan *GroupCastReq) {
	for {
		arg := <-c
		if group := b.Group(arg.gid); group != nil {
			group.Push(arg.proto)
		}
	}
}
