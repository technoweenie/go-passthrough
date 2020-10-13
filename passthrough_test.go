package passthrough

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPassHeaders(t *testing.T) {
	srv := SetupTest()
	defer srv.Close()

	res, w := Get(srv)

	p := New([]string{"Content-Type"})
	p.PassHeaders(res.Header, w)

	head := w.Header()
	if value := head.Get("Content-Type"); value != "text/plain" {
		t.Errorf("Invalid Content Type: %s", value)
	}

	if _, ok := head["Content-Length"]; ok {
		t.Errorf("Passed Content-Length through: %s", head.Get("Content-Length"))
	}

	if w.Status > 0 {
		t.Errorf("How did the status get set?  %d", w.Status)
	}
}

func TestPass(t *testing.T) {
	srv := SetupTest()
	defer srv.Close()

	res, w := Get(srv)

	p := New([]string{"Content-Type"})
	p.Pass(res, w, 200)

	if w.Status != 200 {
		t.Errorf("Invalid status: %d", w.Status)
	}

	head := w.Header()
	if value := head.Get("Content-Type"); value != "text/plain" {
		t.Errorf("Invalid Content Type: %s", value)
	}

	if value := head.Get("Content-Length"); value != "2" {
		t.Errorf("Invalid Content Length: %s", value)
	}

	if body := w.BodyString(); body != "ok" {
		t.Errorf("Invalid body: %s", body)
	}

	if _, ok := head["ETag"]; ok {
		t.Errorf("Passed ETag through: %s", head.Get("ETag"))
	}
}

func TestPassWithTransferEncoding(t *testing.T) {
	srv := SetupTest()
	defer srv.Close()

	res, w := GetWithTransferEncoding(srv)

	p := New([]string{"Content-Type"})
	p.Pass(res, w, 200)

	if w.Status != 200 {
		t.Errorf("Invalid status: %d", w.Status)
	}

	head := w.Header()
	if value := head.Get("Content-Type"); value != "text/plain" {
		t.Errorf("Invalid Content Type: %s", value)
	}

	if value := head.Get("Content-Length"); value != "" {
		t.Errorf("Passed Content-Length through: %s", value)
	}

	if value := head.Get("Transfer-Encoding"); value != "chunked" {
		t.Errorf("Invalid Transfer Encoding: %s", value)
	}

	if body := w.BodyString(); body != "ok" {
		t.Errorf("Invalid body: %s", body)
	}

	if _, ok := head["ETag"]; ok {
		t.Errorf("Passed ETag through: %s", head.Get("ETag"))
	}
}

func SetupTest() *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/example", exampleRequest)
	mux.HandleFunc("/example2", exampleRequest2)
	return httptest.NewServer(mux)
}

func exampleRequest(w http.ResponseWriter, r *http.Request) {
	head := w.Header()
	head.Set("ETag", `"abc"`)
	head.Set("Content-Type", "text/plain")
	head.Set("Content-Length", "2")
	w.WriteHeader(200)
	w.Write([]byte("ok"))
}

func exampleRequest2(w http.ResponseWriter, r *http.Request) {
	head := w.Header()
	head.Set("ETag", `"abc"`)
	head.Set("Content-Type", "text/plain")
	head.Set("Content-Length", "-1")
	w.WriteHeader(200)
	w.Write([]byte("ok"))
}

func Get(srv *httptest.Server) (*http.Response, *FakeResponseWriter) {
	res, err := http.Get(srv.URL + "/example")
	if err != nil {
		panic(err)
	}

	return res, &FakeResponseWriter{new(bytes.Buffer), make(http.Header), 0}
}

func GetWithTransferEncoding(srv *httptest.Server) (*http.Response, *FakeResponseWriter) {
	res, err := http.Get(srv.URL + "/example2")
	if err != nil {
		panic(err)
	}

	return res, &FakeResponseWriter{new(bytes.Buffer), make(http.Header), 0}
}

type FakeResponseWriter struct {
	Buffer *bytes.Buffer
	header http.Header
	Status int
}

func (w *FakeResponseWriter) Header() http.Header {
	return w.header
}

func (w *FakeResponseWriter) WriteHeader(status int) {
	w.Status = status
}

func (w *FakeResponseWriter) Write(buf []byte) (int, error) {
	return w.Buffer.Write(buf)
}

func (w *FakeResponseWriter) BodyString() string {
	return w.Buffer.String()
}
