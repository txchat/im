package tools

import (
	"encoding/json"
	"time"
)

type HttpParams struct {
	Method    string            `json:"Method"`
	ReqUrl    string            `json:"ReqUrl"`
	StrParams string            `json:"StrParams"`
	HeaderMap map[string]string `json:"HeaderMap"`
	Timeout   time.Duration     `json:"Timeout"`
}

func HttpParamsUnmarshal(data []byte) (*HttpParams, error) {
	tHttpParams := &HttpParams{}
	err := json.Unmarshal(data, tHttpParams)
	if err != nil {
		return nil, err
	}

	return tHttpParams, nil
}