package gateway

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/zhangyongxianggithub/grpc-relay/relay/config"
	"golang.org/x/xerrors"
)

type Forwarder struct {
	client  *http.Client
	gateway *config.Gateway
}

func NewForwarder(gateway *config.Gateway) *Forwarder {
	return &Forwarder{
		client:  &http.Client{},
		gateway: gateway,
	}
}

func (forwarder *Forwarder) Forward(writer http.ResponseWriter, request *http.Request) {
	requestId := request.URL.Query().Get("Request-Id")
	if requestId == "" {
		requestId = uuid.New().String()
	}
	startTime := time.Now()
	defer func() {
		duration := time.Now().Sub(startTime)
		slog.Info("eb request complete.", "path", request.URL.Path, "requestId", requestId, "time",
			fmt.Sprintf("%dms", duration.Milliseconds()))
		if err := recover(); err != nil {
			slog.Error("eb request failed.", "path", request.URL.Path, "requestId", requestId, "err", err)
		}
	}()
	forwardUrl := &url.URL{
		Scheme:      forwarder.gateway.Schema,
		Opaque:      request.URL.Opaque,
		User:        request.URL.User,
		Host:        forwarder.gateway.GetEndpoint(),
		Path:        request.URL.Path,
		RawPath:     request.URL.RawPath,
		OmitHost:    request.URL.OmitHost,
		ForceQuery:  request.URL.ForceQuery,
		RawQuery:    request.URL.RawQuery,
		Fragment:    request.URL.Fragment,
		RawFragment: request.URL.RawFragment,
	}
	forwardRequest := &http.Request{
		Method:           request.Method,
		URL:              forwardUrl,
		Proto:            request.Proto,
		ProtoMajor:       request.ProtoMajor,
		ProtoMinor:       request.ProtoMinor,
		Header:           request.Header,
		Body:             request.Body,
		GetBody:          request.GetBody,
		ContentLength:    request.ContentLength,
		TransferEncoding: request.TransferEncoding,
		Close:            request.Close,
		Host:             forwarder.gateway.GetEndpoint(),
		Form:             request.Form,
		PostForm:         request.PostForm,
		MultipartForm:    request.MultipartForm,
		Trailer:          request.Trailer,
		RemoteAddr:       request.RemoteAddr,
		TLS:              request.TLS,
		Response:         request.Response,
	}
	response, err := forwarder.client.Do(forwardRequest)
	if err != nil {
		handleErrorForwardRequest(request, writer, err)
		return
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)
	responseTime := time.Now()
	defer func() {
		duration := time.Now().Sub(responseTime)
		slog.Info(" request response complete.", "requestId", requestId,
			"start time", responseTime.String(), "time",
			fmt.Sprintf("%dms", duration.Milliseconds()))
		if err := recover(); err != nil {
			slog.Error("request failed.", "path", request.URL.Path, "requestId", requestId, "err", err)
		}
	}()
	for key, values := range response.Header {
		writer.Header().Set(key, strings.Join(values, ","))
	}
	writer.WriteHeader(response.StatusCode)
	bytesRead := 0
	buf := make([]byte, 1024)
	flusher, _ := writer.(http.Flusher)
	for {
		size, err := response.Body.Read(buf)
		bytesRead += size
		if size > 0 {
			_, _ = writer.Write(buf[:size])
			if flusher != nil {
				flusher.Flush()
			}
		}
		if err != nil {
			if err == io.EOF {
				break
			} else {
				slog.Error("read eb response failed", slog.Any("error", err))
				break
			}
		}
	}
}

func handleErrorForwardRequest(request *http.Request, writer http.ResponseWriter, err error) {
	xerr := xerrors.New(err.Error())
	slog.Error("forward eb request failed", slog.Any("error", xerr))
	writer.WriteHeader(http.StatusInternalServerError)
	writer.Header().Set("Content-Type", "application/json")
	resp := &ErrorResponse{
		Error:   xerr.Error(),
		Path:    request.URL.Path,
		Headers: request.Header,
	}
	respBytes, err := json.Marshal(resp)
	if err != nil {
		slog.Error("serialize forward eb request failed", slog.Any("error", err))
		return
	}
	_, err = writer.Write(respBytes)
	if err != nil {
		slog.Error("write error response failed", slog.Any("error", err))
		return
	}
}

type ErrorResponse struct {
	Error   string              `json:"error"`
	Path    string              `json:"path"`
	Headers map[string][]string `json:"headers"`
}
