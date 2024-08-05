package eb

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/xerrors"
)

type Forwarder struct {
	client *http.Client
}

func NewForwarder() *Forwarder {
	return &Forwarder{
		client: &http.Client{},
	}
}

func (forwarder *Forwarder) Forward(writer http.ResponseWriter, request *http.Request) {
	forwardUrl := &url.URL{
		Scheme:      "https",
		Opaque:      request.URL.Opaque,
		User:        request.URL.User,
		Host:        "aip.baidubce.com",
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
		Host:             "aip.baidubce.com",
		Form:             request.Form,
		PostForm:         request.PostForm,
		MultipartForm:    request.MultipartForm,
		Trailer:          request.Trailer,
		RemoteAddr:       request.RemoteAddr,
		TLS:              request.TLS,
		Cancel:           request.Cancel,
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
