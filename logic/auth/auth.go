package auth

import "time"

var execAuth = make(map[string]CreateFunc)

type CreateFunc func(url string, timeout time.Duration) Auth

func Register(name string, exec CreateFunc) {
	execAuth[name] = exec
}

func Load(name string) (CreateFunc, error) {
	exec, ok := execAuth[name]
	if !ok {
		return nil, nil
	}
	return exec, nil
}

type Auth interface {
	DoAuth(token string, ext []byte) (string, string, error)
}
