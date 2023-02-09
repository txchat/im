package acc

import (
	"encoding/hex"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/txchat/imparse/proto/signal"

	"github.com/golang/protobuf/proto"
	authProto "github.com/txchat/imparse/proto/auth"
	"github.com/txchat/pkg/auth"
)

func Test_talkClient_DoAuthReConnect(t *testing.T) {
	a := &talkClient{
		url:     "http://172.16.101.107:8888/user/login",
		timeout: time.Second * 5,
	}
	privKey, err := hex.DecodeString("095d1bc1ab3be2047f55b0bef21a8279cc1c09b74aa5853d7c9dbaca4cb99735")
	if err != nil {
		t.Error(err)
		return
	}
	pubKey, err := hex.DecodeString("02ba2c7c208d644c81bd86fb3476adbadacede7717a12e6dade54f3b5d2fcca7ab")
	if err != nil {
		t.Error(err)
		return
	}

	authenticator := auth.NewDefaultAPIAuthenticator()
	token := authenticator.Request("", pubKey, privKey)

	devInfo := &authProto.Login{
		Device:      authProto.Device_Android,
		Username:    "测试comet",
		DeviceToken: "",
		ConnType:    authProto.Login_Reconnect,
		Uuid:        "7447ecc2-948b-465a-bd8a-4830da2e5a09",
		DeviceName:  "虚拟驱动测试",
	}
	extData, _ := proto.Marshal(devInfo)
	gotUid, err := a.DoAuth(token, extData)
	t.Log(err)
	t.Log(gotUid)
}

func Test_talkClient_DoAuthReConnect2(t *testing.T) {
	a := &talkClient{
		url:     "http://172.16.101.107:8888/user/login",
		timeout: time.Second * 5,
	}
	privKey, err := hex.DecodeString("095d1bc1ab3be2047f55b0bef21a8279cc1c09b74aa5853d7c9dbaca4cb99735")
	if err != nil {
		t.Error(err)
		return
	}
	pubKey, err := hex.DecodeString("02ba2c7c208d644c81bd86fb3476adbadacede7717a12e6dade54f3b5d2fcca7ab")
	if err != nil {
		t.Error(err)
		return
	}

	authenticator := auth.NewDefaultAPIAuthenticator()
	token := authenticator.Request("", pubKey, privKey)

	gotUid, err := a.DoAuth(token, nil)
	t.Log(err)
	t.Log(gotUid)
}

func Test_talkClient_DoAuthConnect(t *testing.T) {
	a := &talkClient{
		url:     "http://172.16.101.107:8888/user/login",
		timeout: time.Second * 5,
	}
	privKey, err := hex.DecodeString("095d1bc1ab3be2047f55b0bef21a8279cc1c09b74aa5853d7c9dbaca4cb99735")
	if err != nil {
		t.Error(err)
		return
	}
	pubKey, err := hex.DecodeString("02ba2c7c208d644c81bd86fb3476adbadacede7717a12e6dade54f3b5d2fcca7ab")
	if err != nil {
		t.Error(err)
		return
	}

	authenticator := auth.NewDefaultAPIAuthenticator()
	token := authenticator.Request("", pubKey, privKey)

	devInfo := &authProto.Login{
		Device:      authProto.Device_Android,
		Username:    "1P2vRmRcxNwgSyef12cZxJHqk7sv873tL7",
		DeviceToken: "",
		ConnType:    authProto.Login_Reconnect,
		Uuid:        "3ade6a21-a0d7-48ce-94a2-2f3567adc468",
		DeviceName:  "虚拟驱动",
	}
	extData, _ := proto.Marshal(devInfo)
	gotUid, err := a.DoAuth(token, extData)
	t.Log(err)
	t.Log(gotUid)
}

func TestErrorReject_Encoding(t *testing.T) {
	s := &signal.SignalEndpointLogin{
		Datetime:   0,
		Device:     0,
		DeviceName: "test device",
		Uuid:       "123456789",
	}
	respData := &AuthErrorDataReconnectNotAllowed{
		Code: 0,
		Message: struct {
			Datetime   int64  `json:"datetime"`
			Device     int    `json:"device"`
			DeviceName string `json:"deviceName"`
			Uuid       string `json:"uuid"`
		}{
			Datetime:   0,
			Device:     0,
			DeviceName: "test device",
			Uuid:       "123456789",
		},
		Service: "",
	}
	e := errorDetail(respData)
	srcStr, err := e.Encoding()
	assert.Nil(t, err)
	decoded, err := DecodingErrorReject(srcStr)
	assert.Nil(t, err)
	data, err := proto.Marshal(s)
	assert.Nil(t, err)
	assert.Equal(t, data, decoded)
}
