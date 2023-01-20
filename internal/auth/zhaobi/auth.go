package acc

import (
	"encoding/json"
	"errors"
	"time"

	tools2 "github.com/txchat/im/internel/auth/tools"
)

type talkClient struct {
	url     string
	timeout time.Duration
}

func (a *talkClient) DoAuth(token string, ext []byte) (uid, errMsg string, err error) {
	var (
		bytes []byte
	)
	headers := map[string]string{}
	headers["Authorization"] = token
	bytes, err = tools2.HttpReq(&tools2.HttpParams{
		Method:    "GET",
		ReqUrl:    a.url,
		HeaderMap: headers,
		Timeout:   a.timeout,
	})
	if err != nil {
		return
	}

	var res map[string]interface{}
	err = json.Unmarshal(bytes, &res)
	if err != nil {
		return
	}

	if e, ok := res["error"]; ok {
		err = errors.New(e.(string))
		return
	}

	if _, ok := res["data"]; !ok {
		err = errors.New("invalid auth res")
		return
	}

	data, ok := res["data"].(map[string]interface{})
	if !ok {
		err = errors.New("invalid auth data format")
		return
	}

	if _, ok := data["user_id"]; !ok {
		err = errors.New("invalid auth data")
		return
	}

	uid, ok = data["user_id"].(string)
	if !ok {
		err = errors.New("invalid auth data id format")
		return
	}

	return
}
