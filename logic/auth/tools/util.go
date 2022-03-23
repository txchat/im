package tools

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

var client = &http.Client{}

func HttpReq(tHttpParams *HttpParams) ([]byte, error) {
	var req *http.Request
	var err error

	switch tHttpParams.Method {
	case HTTP_GET, HTTP_DELETE:
		if "" == tHttpParams.StrParams {
			req, err = http.NewRequest(tHttpParams.Method, tHttpParams.ReqUrl, nil)
		} else {
			req, err = http.NewRequest(tHttpParams.Method, tHttpParams.ReqUrl+"?"+tHttpParams.StrParams, nil)
		}
	case HTTP_POST:
		req, err = http.NewRequest(tHttpParams.Method, tHttpParams.ReqUrl, strings.NewReader(tHttpParams.StrParams))
	default:
		panic(tHttpParams.Method)
	}
	if nil != err {
		return nil, err
	}

	for k, v := range tHttpParams.HeaderMap {
		req.Header.Add(k, v)
	}

	if 0 != int64(tHttpParams.Timeout) {
		client.Timeout = tHttpParams.Timeout
	} else {
		client.Timeout = HTTP_DEFAULT_TIMEOUT
	}

	resp, err := client.Do(req)
	if nil != resp {
		defer resp.Body.Close()
	}

	if nil != err {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(fmt.Sprintf("http request errcode: %v", resp.StatusCode))
	}

	return ioutil.ReadAll(resp.Body)
}
