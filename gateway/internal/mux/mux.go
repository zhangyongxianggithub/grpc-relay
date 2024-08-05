package mux

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

type Middleware func(http.Handler) http.Handler

type ChainMux struct {
	*http.ServeMux
	middlewares []Middleware
}

func NewChainMux() *ChainMux {
	return &ChainMux{
		ServeMux: http.NewServeMux(),
	}
}

func (chainMux *ChainMux) Use(middlewares ...Middleware) {
	chainMux.middlewares = append(chainMux.middlewares, middlewares...)
}

func (chainMux *ChainMux) ApplyMiddlewares(handler http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

func (chainMux *ChainMux) Handle(pattern string, handler http.Handler) {
	chainHandler := chainMux.ApplyMiddlewares(handler, chainMux.middlewares...)
	chainMux.ServeMux.Handle(pattern, chainHandler)
}

func (chainMux *ChainMux) HandleFunc(pattern string, handler http.HandlerFunc) {
	chainHandler := chainMux.ApplyMiddlewares(handler, chainMux.middlewares...)
	chainMux.ServeMux.Handle(pattern, chainHandler)
}

func (chainMux *ChainMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	chainMux.ServeMux.ServeHTTP(w, r)
}

type CachedRequest struct {
	r    *http.Request
	body []byte
}

func NewCachedRequest(r *http.Request) (*CachedRequest, error) {
	buf, _ := io.ReadAll(r.Body)
	newBody := io.NopCloser(bytes.NewBuffer(buf))
	bodyReplica := io.NopCloser(bytes.NewBuffer(buf))
	r.Body = newBody
	body, err := io.ReadAll(bodyReplica)
	if err != nil {
		return nil, err
	}
	if body == nil {
		body = []byte{}
	}
	return &CachedRequest{
		r: &http.Request{
			Method:           r.Method,
			URL:              r.URL,
			Proto:            r.Proto,
			ProtoMajor:       r.ProtoMajor,
			ProtoMinor:       r.ProtoMinor,
			Header:           r.Header,
			Body:             io.NopCloser(bytes.NewBuffer(buf)),
			GetBody:          r.GetBody,
			ContentLength:    r.ContentLength,
			TransferEncoding: r.TransferEncoding,
			Close:            r.Close,
			Host:             r.Host,
			Form:             r.Form,
			PostForm:         r.PostForm,
			MultipartForm:    r.MultipartForm,
			Trailer:          r.Trailer,
			RemoteAddr:       r.RemoteAddr,
			RequestURI:       r.RequestURI,
			TLS:              r.TLS,
			Response:         r.Response,
		},
		body: body,
	}, nil
}

func (r *CachedRequest) JSON() (map[string]any, error) {
	m := map[string]any{}
	err := json.Unmarshal(r.body, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (r *CachedRequest) String() string {
	return string(r.body)
}

func (r *CachedRequest) Path() string {
	return r.r.URL.Path
}

func (r *CachedRequest) Header(header string) string {
	return r.r.Header.Get(header)
}

type CachedResponse struct {
	Request    *CachedRequest
	Writer     http.ResponseWriter
	StatusCode int
	Response   []byte
}

func (c *CachedResponse) Header() http.Header {
	return c.Writer.Header()
}

func (c *CachedResponse) Write(bytes []byte) (int, error) {
	c.Response = append(c.Response, bytes...)
	return c.Writer.Write(bytes)
}

func (c *CachedResponse) WriteHeader(statusCode int) {
	c.Writer.WriteHeader(statusCode)
	c.StatusCode = statusCode
}

func (c *CachedResponse) Flush() {
	if flusher, ok := c.Writer.(http.Flusher); ok {
		flusher.Flush()
	}
}

func RecoveryInterceptor(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				slog.Error("internal server error", "path", r.URL.Path, slog.Any("error", err))
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				serverError := &InternalServerError{Error: fmt.Sprintf("%v", err)}
				errorBody, _ := json.Marshal(serverError)
				_, _ = w.Write(errorBody)
			}
		}()
		handler.ServeHTTP(w, r)
	})
}

type InternalServerError struct {
	Error string `json:"error"`
}
