package tools

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

var client = &http.Client{}

func HTTPReq(tHTTPParams *HTTPParams) ([]byte, error) {
	var req *http.Request
	var err error

	switch tHTTPParams.Method {
	case HTTPGet, HTTPDelete:
		if "" == tHTTPParams.StrParams {
			req, err = http.NewRequest(tHTTPParams.Method, tHTTPParams.ReqURL, nil)
		} else {
			req, err = http.NewRequest(tHTTPParams.Method, tHTTPParams.ReqURL+"?"+tHTTPParams.StrParams, nil)
		}
	case HTTPPost:
		req, err = http.NewRequest(tHTTPParams.Method, tHTTPParams.ReqURL, strings.NewReader(tHTTPParams.StrParams))
	default:
		panic(tHTTPParams.Method)
	}
	if nil != err {
		return nil, err
	}

	for k, v := range tHTTPParams.HeaderMap {
		req.Header.Add(k, v)
	}

	if 0 != int64(tHTTPParams.Timeout) {
		client.Timeout = tHTTPParams.Timeout
	} else {
		client.Timeout = HTTPDefaultTimeout
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
