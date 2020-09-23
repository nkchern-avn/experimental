package httpext

import (
	"encoding/json"
	"log"
	"net/http"
)

var (
	Log          Logger       = &defaultLogger{}
	errorHandler ErrorHandler = defautErrorHandler

	builder HandlerBuilder = &defaultBuilder{}
)

type Response interface {
	Status() int
	Bytes() []byte
}

type HandlerBuilder interface {
	Handle(JSONHandler) func(http.ResponseWriter, *http.Request)
}

type defaultBuilder struct{}

func (b *defaultBuilder) Handle(fn JSONHandler) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, DELETE, PUT")
			w.Header().Set("Access-Control-Allow-Headers", "*")
			w.WriteHeader(http.StatusOK)

			return
		}

		resp, err := fn(r)
		// unhandled error
		if err != nil {
			respondInternalError(r, w, err)
			return
		}

		w.WriteHeader(resp.Status())
		if _, err := w.Write(resp.Bytes()); err != nil {
			// unable to write response, last attempt to log it
			Log.Error(r, err)
		}
	}
}

type errorHandlingBuilder struct {
	onError ErrorHandler //= defautErrorHandler
}

func (b *errorHandlingBuilder) Handle(fn JSONHandler) func(http.ResponseWriter, *http.Request) {
	return Handle(func(r *http.Request) (Response, error) {
		resp, err := fn(r)
		if err != nil {
			resp, err = b.onError(r, err)
			if err != nil {
				return nil, err
			}
		}
		return resp, nil
	})
}

func WithErroHandler(h ErrorHandler) HandlerBuilder {
	if h == nil {
		h = defautErrorHandler
	}
	return &errorHandlingBuilder{onError: h}
}

type jsonResponse struct {
	status int
	obj    interface{}
}

func (r *jsonResponse) Status() int { return r.status }

func (r *jsonResponse) Bytes() []byte {
	bytes, err := json.Marshal(r.obj)
	if err != nil {
		panic(err)
	}
	return bytes
}

func JSON(o interface{}) Response { return &jsonResponse{obj: o, status: http.StatusOK} }

type textResponse struct {
	status int
	text   string
}

func (r *textResponse) Status() int { return r.status }

func (r *textResponse) Bytes() []byte { return []byte(r.text) }

func Plain(s string) Response { return &textResponse{text: s, status: http.StatusOK} }

type JSONHandler func(*http.Request) (Response, error)

type ErrorHandler func(r *http.Request, err error) (Response, error)

type Logger interface {
	Error(r *http.Request, err error)
}

func Handle(fn JSONHandler) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, DELETE, PUT")
			w.Header().Set("Access-Control-Allow-Headers", "*")
			w.WriteHeader(http.StatusOK)

			return
		}

		resp, err := fn(r)

		if err != nil {
			resp, err = errorHandler(r, err)

			// unhandled error
			if err != nil {
				respondInternalError(r, w, err)
				return
			}
		}

		w.WriteHeader(resp.Status())
		if _, err := w.Write(resp.Bytes()); err != nil {
			// unable to write response, last attempt to log it
			Log.Error(r, err)
		}
	}
}

func respondInternalError(r *http.Request, w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	Log.Error(r, err)
}

type defaultLogger struct{}

func (e *defaultLogger) Error(r *http.Request, err error) {
	log.Printf("ERROR %s %s", r.URL, err)
}

func defautErrorHandler(r *http.Request, err error) (Response, error) {
	return nil, err
	//	return &jsonResponse{obj: map[string]string{}, status: http.StatusInternalServerError}, nil
}
