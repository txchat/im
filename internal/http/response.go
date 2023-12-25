package websocket

import (
	"bufio"
	"net"
	"net/http"
)

type response struct {
	rw   *bufio.ReadWriter
	req  *http.Request
	conn net.Conn

	body       []byte
	statusCode int
	header     http.Header
}

func (w *response) Header() http.Header {
	return w.header
}

func (w *response) Write(b []byte) (int, error) {
	w.body = b
	return 0, nil
}

func (w *response) WriteHeader(statusCode int) {
	w.statusCode = statusCode
}

func (w *response) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.conn, w.rw, nil
}

func ReadRequest(rw *bufio.ReadWriter, req *http.Request, conn net.Conn) (*response, error) {
	w := &response{
		rw:         rw,
		req:        req,
		conn:       conn,
		body:       []byte{},
		statusCode: 0,
		header:     map[string][]string{},
	}
	return w, nil
}
