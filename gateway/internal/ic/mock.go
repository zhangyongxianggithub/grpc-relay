package ic

import (
	"context"
	"strconv"

	"github.com/zhangyongxianggithub/grpc-relay/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Inference struct {
	pb.UnimplementedGRPCInferenceServiceServer
}

func (inference *Inference) ModelStreamInfer(modelInferRequest *pb.ModelInferRequest,
	modelStreamInferServer pb.GRPCInferenceService_ModelStreamInferServer) error {

	for i := 0; i < 10; i++ {
		responseSegment := &pb.ModelInferResponse{
			RequestId:  modelInferRequest.RequestId,
			SentenceId: int32(i),
			Output:     modelInferRequest.Input + ": " + strconv.Itoa(i),
			ModelId:    modelInferRequest.ModelId,
			TraceId:    modelInferRequest.TraceId,
			TenantId:   modelInferRequest.TenantId,
			OutputType: modelInferRequest.InputType,
		}
		if err := modelStreamInferServer.Send(responseSegment); err != nil {
			return err
		}
	}
	return nil
}
func (inference *Inference) ModelFetchRequest(ctx context.Context,
	modelFetchRequestParams *pb.ModelFetchRequestParams) (*pb.ModelFetchRequestResult, error) {
	result := new(pb.ModelFetchRequestResult)
	result.Requests = make([]*pb.ModelInferRequest, 0)
	for _, model := range modelFetchRequestParams.ModelId {
		result.Requests = append(result.Requests, &pb.ModelInferRequest{
			ModelId:   model,
			RequestId: "zhangyongxiang",
			TraceId:   "zhangyongxiang",
			Input:     "zhangyongxiang",
			TenantId:  "zhangyongxiang",
			InputType: 0,
		})
	}
	return result, nil
}
func (inference *Inference) ModelSendResponse(modelSendResponseServer pb.GRPCInferenceService_ModelSendResponseServer) error {
	return status.Errorf(codes.Unimplemented, "method ModelSendResponse not implemented")
}
