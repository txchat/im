package websocket

import "bufio"

func peek(w *bufio.Writer, n int) ([]byte, error) {
	var err error
	if n < 0 {
		return nil, bufio.ErrNegativeCount
	}
	if n > w.Size() {
		return nil, bufio.ErrBufferFull
	}
	for w.Available() < n && err == nil {
		err = w.Flush()
	}
	if err != nil {
		return nil, err
	}
	buf := w.AvailableBuffer()
	return buf[:n], nil
}

func pop(r *bufio.Reader, n int) ([]byte, error) {
	buf, err := r.Peek(n)
	if err != nil {
		return nil, err
	}
	if _, err = r.Discard(n); err != nil {
		return nil, err
	}
	return buf, nil
}
