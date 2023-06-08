package comet

import (
	"fmt"
	"sync"

	"github.com/txchat/im/api/protocol"
)

// Group is a group and store channel group info.
type Group struct {
	ID        string
	rLock     sync.RWMutex
	next      *Node
	drop      bool
	Online    int32 // dirty read is ok
	AllOnline int32
}

// NewGroup new a group struct, store channel group info.
func NewGroup(id string) (r *Group) {
	r = new(Group)
	r.ID = id
	r.drop = false
	r.next = nil
	r.Online = 0
	return
}

// Put hold the channel instance.
func (r *Group) Put(ch *Channel) (err error) {
	r.rLock.Lock()
	if !r.drop && ch.GetNode(r.ID) == nil {
		node := &Node{
			Current: ch,
			Next:    nil,
			Prev:    nil,
		}
		ch.SetNode(r.ID, node)
		if r.next != nil {
			r.next.Prev = node
		}
		node.Next = r.next
		node.Prev = nil
		r.next = node // insert to header
		r.Online++
	} /* else { 	//del: 2021年7月21日16:32:54 dld
		err = errors.ErrGroupDroped
	}*/
	r.rLock.Unlock()
	return
}

// Del delete channel from the group.
func (r *Group) Del(ch *Channel) bool {
	r.rLock.Lock()
	if node := ch.GetNode(r.ID); node != nil {
		if node.Next != nil {
			// if not footer
			node.Next.Prev = node.Prev
		}
		if node.Prev != nil {
			// if not header
			node.Prev.Next = node.Next
		} else {
			r.next = node.Next
		}
		r.Online--
		//r.drop = (r.Online == 0)
		ch.DelNode(r.ID) //2021年6月10日 删除对应node，防止再次put的时候报ErrGroupDroped错误；dld
	}
	r.rLock.Unlock()
	return r.drop
}

// Push broadcast msg inner group, if chan full discard it.
func (r *Group) Push(p *protocol.Proto) {
	r.rLock.RLock()
	for node := r.next; node != nil; node = node.Next {
		if node.Current != nil {
			_ = node.Current.Push(p)
		}
	}
	r.rLock.RUnlock()
}

// Members group members Key,IP
func (r *Group) Members() ([]string, []string) {
	r.rLock.RLock()
	members := make([]string, 0)
	mIp := make([]string, 0)
	for node := r.next; node != nil; node = node.Next {
		if node.Current != nil {
			members = append(members, node.Current.Key)
			mIp = append(mIp, fmt.Sprintf("%v:%v", node.Current.IP, node.Current.Port))
		}
	}
	r.rLock.RUnlock()
	return members, mIp
}

// Close remove group channel index.
func (r *Group) Close() {
	r.rLock.Lock()
	for node := r.next; node != nil; node = node.Next {
		if ch := node.Current; ch != nil {
			ch.DelNode(r.ID)
		}
	}
	r.rLock.Unlock()
}

// OnlineNum the group all online.
func (r *Group) OnlineNum() int32 {
	if r.AllOnline > 0 {
		return r.AllOnline
	}
	return r.Online
}
