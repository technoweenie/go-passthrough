package passthrough

import (
	"io"
	"net/http"
	"strconv"
)

type HttpPassthrough struct {
	Headers []string
}

func New(headers []string) *HttpPassthrough {
	return &HttpPassthrough{headers}
}

func (p *HttpPassthrough) Pass(res *http.Response, w http.ResponseWriter, status int) {
	head := p.PassHeaders(res.Header, w)
	head.Set("Content-Length", strconv.Itoa(int(res.ContentLength)))
	w.WriteHeader(status)
	io.Copy(w, res.Body)
	res.Body.Close()
}

func (p *HttpPassthrough) PassHeaders(resHeader http.Header, w http.ResponseWriter) http.Header {
	head := w.Header()
	for _, header := range p.Headers {
		if value := resHeader.Get(header); len(value) != 0 {
			head.Set(header, value)
		}
	}
	return head
}
