package passthrough

import (
	"io"
	"net/http"
	"strconv"
)

type Client struct {
	Headers []string
}

func New(headers ...string) *Client {
	return &Client{headers}
}

func (c *Client) Pass(res *http.Response, w http.ResponseWriter, status int) {
	head := c.PassHeaders(res.Header, w)
	head.Set("Content-Length", strconv.Itoa(int(res.ContentLength)))
	w.WriteHeader(status)
	io.Copy(w, res.Body)
	res.Body.Close()
}

func (c *Client) PassHeaders(resHeader http.Header, w http.ResponseWriter) http.Header {
	head := w.Header()
	for _, header := range c.Headers {
		if value := resHeader.Get(header); len(value) != 0 {
			head.Set(header, value)
		}
	}
	return head
}
