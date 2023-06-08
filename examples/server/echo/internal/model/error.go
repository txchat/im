package model

import (
	"errors"
)

var (
	ErrAppID            = errors.New("app_id not compared")
	ErrCustomNotSupport = errors.New("biz not support")
)
