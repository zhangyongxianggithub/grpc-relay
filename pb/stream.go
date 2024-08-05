package pb

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/protobuf/proto"
)

func ForwardResponseStream(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler,
	w http.ResponseWriter, req *http.Request, recv func() (proto.Message, error),
	opts ...func(context.Context, http.ResponseWriter, proto.Message) error) {
	md, ok := runtime.ServerMetadataFromContext(ctx)
	if !ok {
		grpclog.Error("Failed to extract ServerMetadata from context")
		http.Error(w, "unexpected error", http.StatusInternalServerError)
		return
	}
	handleForwardResponseServerMetadata(w, mux, md)

	// w.Header().Set("Transfer-Encoding", "chunked")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	flusher, ok := w.(http.Flusher)
	if !ok {
		handleForwardResponseStreamError(ctx, marshaler, w, req, mux, errors.New("Streaming unsupported!"))
		return
	}
	for {
		resp, err := recv()
		if errors.Is(err, io.EOF) {
			return
		}
		if err != nil {
			handleForwardResponseStreamError(ctx, marshaler, w, req, mux, err)
			return
		}
		if err := handleForwardResponseOptions(ctx, w, resp, opts); err != nil {
			handleForwardResponseStreamError(ctx, marshaler, w, req, mux, err)
			return
		}
		var buf []byte
		buf, err = marshaler.Marshal(resp)
		_, _ = fmt.Fprintf(w, "data: %s\n\n", string(buf))
		flusher.Flush()
	}
}

func handleForwardResponseServerMetadata(w http.ResponseWriter, mux *runtime.ServeMux, md runtime.ServerMetadata) {
	for k, vs := range md.HeaderMD {
		for _, v := range vs {
			w.Header().Add(k, v)
		}
	}
}

func handleForwardResponseOptions(ctx context.Context, w http.ResponseWriter, resp proto.Message,
	opts []func(context.Context, http.ResponseWriter, proto.Message) error) error {
	if len(opts) == 0 {
		return nil
	}
	for _, opt := range opts {
		if err := opt(ctx, w, resp); err != nil {
			grpclog.Errorf("Error handling ForwardResponseOptions: %v", err)
			return err
		}
	}
	return nil
}

func handleForwardResponseStreamError(ctx context.Context, marshaler runtime.Marshaler,
	w http.ResponseWriter, req *http.Request, mux *runtime.ServeMux, err error) {
	w.WriteHeader(runtime.HTTPStatusFromCode(500))
	buf, err := marshaler.Marshal(err.Error())
	if err != nil {
		grpclog.Errorf("Failed to marshal an error: %v", err)
		return
	}
	if _, err := w.Write(buf); err != nil {
		grpclog.Errorf("Failed to notify error to client: %v", err)
		return
	}
}
