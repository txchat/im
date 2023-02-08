package comet

import (
	"errors"

	"github.com/txchat/im/api/protocol"
)

var (
	ErrNotGroupMember     = errors.New("not group members")
	ErrNotFriend          = errors.New("target not friend members")
	ErrUnsupportedChannel = errors.New("unsupported channel type")
)

type Filter struct {
	ch *Channel
}

func NewFilter(ch *Channel) *Filter {
	return &Filter{ch: ch}
}

func (f *Filter) Filter(channel protocol.Channel, target string) error {
	switch channel {
	case protocol.Channel_Private:
	case protocol.Channel_Group:
		if g := f.ch.GetNode(target); g == nil {
			return ErrNotGroupMember
		}
	default:
		return ErrUnsupportedChannel
	}
	return nil
}
