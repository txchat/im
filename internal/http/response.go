package websocket

import (
	"bufio"
	"net/http"
)

type response struct {
}

func (w *response) Header() http.Header {
	panic("not implemented") // TODO: Implement
}

func (w *response) Write(_ []byte) (int, error) {
	panic("not implemented") // TODO: Implement
}

func (w *response) WriteHeader(statusCode int) {
	panic("not implemented") // TODO: Implement
}

func ReadRequest(r *bufio.Writer, req *http.Request) (*response, error) {
	rw := &response{}
	return rw, nil
}
