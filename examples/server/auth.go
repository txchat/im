package server

import "errors"

var ErrUserNotExists = errors.New("user not exists")

type mock struct {
}

func (c *mock) DoAuth(token string, ext []byte) (uid string, err error) {
	if token == "" {
		return "", ErrUserNotExists
	}
	return token, nil
}
