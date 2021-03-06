package errors

import (
	"errors"
)

// .
var (
	// server
	ErrHandshake = errors.New("handshake failed")
	ErrOperation = errors.New("request operation not valid")
	// ring
	ErrRingEmpty = errors.New("ring buffer empty")
	ErrRingFull  = errors.New("ring buffer full")
	// timer
	ErrTimerFull   = errors.New("timer full")
	ErrTimerEmpty  = errors.New("timer empty")
	ErrTimerNoItem = errors.New("timer item not exist")
	// channel
	ErrUnconnected  = errors.New("client unconnected error")
	ErrJoinGroupArg = errors.New("rpc joingroup arg error")
	ErrPushMsgArg   = errors.New("rpc pushmsg arg error")
	ErrPushMsgsArg  = errors.New("rpc pushmsgs arg error")
	ErrMPushMsgArg  = errors.New("rpc mpushmsg arg error")
	ErrMPushMsgsArg = errors.New("rpc mpushmsgs arg error")
	// bucket
	ErrBroadCastArg     = errors.New("rpc broadcast arg error")
	ErrBroadCastRoomArg = errors.New("rpc broadcast  room arg error")

	// group
	ErrGroupDroped = errors.New("group droped")
	// rpc
	ErrLogic = errors.New("logic rpc is not available")
)
