package acc

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
	"github.com/txchat/im/logic/auth/tools"
	xproto "github.com/txchat/imparse/proto"
)

var dev = make(map[xproto.Device]string)

type EndpointRejectResp struct {
	Result  int    `json:"result"`
	Message string `json:"message"`
	Data    struct {
		Code    int    `json:"code"`
		Service string `json:"service"`
		Message struct {
			UUid       string `json:"uuid"`
			Device     int    `json:"device"`
			DeviceName string `json:"deviceName"`
			Datetime   int64  `json:"datetime"`
		} `json:"message"`
	} `json:"data"`
}

func errorMetaData(respData *AuthErrorDataReconnectNotAllowed) (string, error) {
	data, err := proto.Marshal(&xproto.SignalEndpointLogin{
		Uuid:       respData.Message.Uuid,
		Device:     xproto.Device(respData.Message.Device),
		DeviceName: respData.Message.DeviceName,
		Datetime:   uint64(respData.Message.Datetime),
	})
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

type talkClient struct {
	url     string
	timeout time.Duration
}

func (a *talkClient) DoAuth(token string, ext []byte) (uid, errMsg string, err error) {
	var (
		bytes     []byte
		strParams string
	)
	headers := map[string]string{}
	headers["FZM-SIGNATURE"] = token

	if len(ext) != 0 {
		var device xproto.Login
		err = proto.Unmarshal(ext, &device)
		if err != nil {
			return "", "", err
		}

		headers["FZM-UUID"] = device.Uuid
		headers["FZM-DEVICE"] = device.Device.String()
		headers["FZM-DEVICE-NAME"] = device.DeviceName
		headers["Content-type"] = "application/json"
		reqBody := gin.H{
			"connType": device.ConnType,
		}
		reqData, err := json.Marshal(reqBody)
		if err != nil {
			return "", "", err
		}
		strParams = string(reqData)
	}

	bytes, err = tools.HttpReq(&tools.HttpParams{
		Method:    "POST",
		ReqUrl:    a.url,
		HeaderMap: headers,
		Timeout:   a.timeout,
		StrParams: strParams,
	})
	if err != nil {
		return
	}

	var res AuthReply
	err = json.Unmarshal(bytes, &res)
	if err != nil {
		return
	}
	log.Debug().Interface("resp", res).Msg("auth reply")

	switch res.Result {
	case 0:
		var success AuthSuccessData
		if err = mapstructure.Decode(res.Data, &success); err != nil {
			err = errors.New("invalid auth success data")
			return
		}
		errMsg = ""
		uid = success.Address
	case -1016:
		var errNotAllowed AuthErrorDataReconnectNotAllowed
		if err = mapstructure.Decode(res.Data, &errNotAllowed); err != nil {
			return
		}
		emsg := ""
		if emsg, err = errorMetaData(&errNotAllowed); err != nil {
			return
		}
		errMsg = emsg
		err = errors.New(res.Message)
	default:
		errMsg = ""
		err = errors.New(res.Message)
	}
	log.Debug().Str("errMsg", errMsg).Interface("err", err).Msg("auth reply code")
	return
}
