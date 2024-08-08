package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/tmaxmax/go-sse"
	"github.com/zhangyongxianggithub/grpc-relay/pb"
	"github.com/zhangyongxianggithub/grpc-relay/relay/config"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Inference struct {
	pb.UnimplementedGRPCInferenceServiceServer
	gateway   *config.Gateway
	marshaler *runtime.HTTPBodyMarshaler
	client    *sse.Client
	prefix    string
}

func (inference *Inference) ModelStreamInfer(modelInferRequest *pb.ModelInferRequest,
	modelStreamInferServer pb.GRPCInferenceService_ModelStreamInferServer) (error error) {
	reqBody, _ := json.Marshal(modelInferRequest)
	slog.Info("relay model infer request", "prefix", inference.prefix, "request", string(reqBody))
	startTime := time.Now()
	defer func() {
		duration := time.Now().Sub(startTime)
		slog.Info("request complete.",
			"requestId", modelInferRequest.RequestId,
			"took time", fmt.Sprintf("%dms", duration.Milliseconds()))
		if err := recover(); err != nil {
			error = fmt.Errorf("%v", err)
		}
	}()
	req, _ := http.NewRequest(http.MethodPost,
		inference.gateway.GetUriPrefix()+
			inference.prefix+"/grpcinferenceservice/modelstreaminfer?token="+inference.gateway.Token,
		bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	conn := inference.client.NewConnection(req)
	count := 0
	responseTime := time.Now()
	defer func() {
		duration := time.Now().Sub(responseTime)
		slog.Info(" request complete.", "start time",
			responseTime.String(), "requestId", modelInferRequest.RequestId,
			"took time", fmt.Sprintf("%dms", duration.Milliseconds()))
		if err := recover(); err != nil {
			error = fmt.Errorf("%v", err)
		}
	}()
	unsubscribe := conn.SubscribeMessages(func(event sse.Event) {
		if count == 0 {
			responseTime = time.Now()
		}
		count++
		slog.Info("receive stream segment.", "requestId", modelInferRequest.RequestId, "count", count)
		response := pb.ModelInferResponse{}
		_ = inference.marshaler.Unmarshal([]byte(event.Data), &response)
		_ = modelStreamInferServer.Send(&response)
	})
	if err := conn.Connect(); !errors.Is(err, io.EOF) {
		return err
	}
	unsubscribe()
	return nil
}
func (inference *Inference) ModelFetchRequest(ctx context.Context,
	modelFetchRequestParams *pb.ModelFetchRequestParams) (result *pb.ModelFetchRequestResult, err error) {
	reqBody, _ := json.Marshal(modelFetchRequestParams)
	slog.Info("relay model fetch request", "prefix", inference.prefix, "request", string(reqBody))
	startTime := time.Now()
	defer func() {
		duration := time.Now().Sub(startTime)
		slog.Info("request complete.",
			"request", modelFetchRequestParams.String(),
			"took time", fmt.Sprintf("%dms", duration.Milliseconds()))
		if err := recover(); err != nil {
			err = fmt.Errorf("%v", err)
		}
	}()
	req, _ := http.NewRequest(http.MethodPost,
		inference.gateway.GetUriPrefix()+
			inference.prefix+"/grpcinferenceservice/modelfetchrequest?token="+inference.gateway.Token,
		bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := inference.client.HTTPClient.Do(req)
	if err != nil {
		slog.Error("relay model fetch request failed", "error", err)
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	body, _ := io.ReadAll(resp.Body)
	response := new(pb.ModelFetchRequestResult)
	err = json.Unmarshal(body, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}
func (inference *Inference) ModelSendResponse(modelSendResponseServer pb.GRPCInferenceService_ModelSendResponseServer) error {
	return status.Errorf(codes.Unimplemented, "method ModelSendResponse not implemented")
}
