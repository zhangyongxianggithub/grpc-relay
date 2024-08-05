package gateway

import (
	"log/slog"
	"net"
	"net/http"

	"bestzyx.com/grpc-relay/pb"
	"bestzyx.com/grpc-relay/relay/config"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/tmaxmax/go-sse"
	"golang.org/x/xerrors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
)

type GrpcServerEngine struct {
	grpcServer *grpc.Server
}

var grpcEngine = new(GrpcServerEngine)

func init() {
	config.AddInitializer(grpcEngine)
	slog.Info("append grpc engine initializer to initializers")
}

func (e *GrpcServerEngine) Name() string {
	return "relay-server"
}

func (e *GrpcServerEngine) Init(config *config.ServerConfig) error {

	marshaler := &runtime.HTTPBodyMarshaler{
		Marshaler: &runtime.JSONPb{MarshalOptions: protojson.MarshalOptions{},
			UnmarshalOptions: protojson.UnmarshalOptions{}},
	}
	client := &sse.Client{
		Backoff: sse.Backoff{
			MaxRetries: -1,
		},
	}
	// create a new grpc server and register the inference gateway
	if config.Server.DataListen != "" {
		go func() {
			lis, err := net.Listen("tcp", config.Server.DataListen)
			if err != nil {
				slog.Error("data server grpc server listen failed", slog.Any("err", err))
				panic(err)
			}
			grpcServer := grpc.NewServer(grpc.Creds(insecure.NewCredentials()))
			pb.RegisterGRPCInferenceServiceServer(grpcServer, &Inference{
				gateway:   config.Gateway,
				marshaler: marshaler,
				client:    client,
				prefix:    "/comate/v1/grpc/dataserver",
			})
			if err = grpcServer.Serve(lis); err != nil {
				slog.Error("data server grpc server start failed", slog.Any("err", err))
				panic(err)
			}
		}()
	}
	if config.Server.EBListen != "" {
		go func() {
			mux := http.NewServeMux()
			forwarder := NewForwarder(config.Gateway)
			mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
				forwarder.Forward(writer, request)
			})
			server := &http.Server{
				Addr:    config.Server.EBListen,
				Handler: mux,
			}
			if err := server.ListenAndServe(); err != nil {
				slog.Error("eb proxy start failed", slog.Any("err", xerrors.New(err.Error())))
				panic(err)
			}
		}()
	}
	lis, err := net.Listen("tcp", config.Server.ICListen)
	if err != nil {
		slog.Error("ic grpc server listen failed", slog.Any("err", err))
		panic(err)
	}
	grpcServer := grpc.NewServer(grpc.Creds(insecure.NewCredentials()))
	pb.RegisterGRPCInferenceServiceServer(grpcServer, &Inference{
		gateway:   config.Gateway,
		marshaler: marshaler,
		client:    client,
		prefix:    "/comate/v1/grpc/ic",
	})
	if err = grpcServer.Serve(lis); err != nil {
		slog.Error("ic grpc server start failed", slog.Any("err", err))
		panic(err)
	}
	return nil
}
