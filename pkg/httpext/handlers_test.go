package httpext

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	port = 38888
)

var (
	host = fmt.Sprintf("http://localhost:%d", port)
)

type response struct {
	X string
	Y int
}

func init() {
	http.HandleFunc("/", Handle(func(r *http.Request) (Response, error) {
		resp := &response{X: "hello", Y: 42}
		return JSON(resp), nil
	}))

	http.HandleFunc("/error", Handle(func(r *http.Request) (Response, error) {
		return nil, errors.New("foo error")
	}))

	go func() {
		err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
		panic(err) // must never happen
	}()

}

func TestShouldHanlde(t *testing.T) {
	resp, err := http.Get(host + "/")

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var actual response
	expected := &response{X: "hello", Y: 42}

	assert.NoError(t, json.NewDecoder(resp.Body).Decode(&actual))
	assert.Equal(t, expected, &actual)
}

func TestShouldHanldeError(t *testing.T) {
	resp, err := http.Get(host + "/error")

	assert.NoError(t, err)
	assert.Equal(t, 500, resp.StatusCode)

	b, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, "", string(b))
}

func TestShouldLogResponseWriterErrors(t *testing.T) {
	l := &testLogger{}
	Log = l
	defer func() { Log = &defaultLogger{} }() // restore log; not good code

	w := &mockWriter{err: errors.New("boom")}

	f := Handle(func(r *http.Request) (Response, error) {
		resp := &response{X: "hello", Y: 42}
		return JSON(resp), nil
	})
	req := &http.Request{}
	f(w, req)

	assert.Equal(t, 200, w.status) // status was set successfully
	assert.Equal(t, w.err, l.logged)
	assert.Equal(t, req, l.r)
}

type mockWriter struct {
	err    error
	status int
}

func (w *mockWriter) Header() http.Header { return http.Header{} }

func (w *mockWriter) Write(_ []byte) (int, error) { return 0, w.err }

func (w *mockWriter) WriteHeader(statusCode int) {
	w.status = statusCode
}

type testLogger struct {
	logged error
	r      *http.Request
}

func (l *testLogger) Error(r *http.Request, err error) {
	l.r = r
	l.logged = err
}
