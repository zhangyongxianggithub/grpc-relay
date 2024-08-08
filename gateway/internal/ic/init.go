package ic

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"net"
	"net/http"
	"time"

	"github.com/esurdam/go-swagger-ui"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/zhangyongxianggithub/grpc-relay/gateway/config"
	"github.com/zhangyongxianggithub/grpc-relay/gateway/internal/auth"
	"github.com/zhangyongxianggithub/grpc-relay/gateway/internal/charge"
	"github.com/zhangyongxianggithub/grpc-relay/gateway/internal/mux"
	"github.com/zhangyongxianggithub/grpc-relay/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcServerAdapter struct {
	grpcServer      *grpc.Server
	ICServers       map[string]*runtime.ServeMux
	DataServers     map[string]*runtime.ServeMux
	ICServerSeeds   []string
	DataServerSeeds []string
	Rand            *rand.Rand
}

var grpcServerAdapter = new(GrpcServerAdapter)

var _ config.Initializer = (*GrpcServerAdapter)(nil)

func init() {
	config.AddInitializer(grpcServerAdapter)
	slog.Info("append grpc-service-adapter to initializers")
}

func (g *GrpcServerAdapter) Name() string {
	return "grpc-service-adapter"
}

func (g *GrpcServerAdapter) Init(config *config.ServerConfig, multiplexer *mux.ChainMux) error {
	g.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	g.DataServers = make(map[string]*runtime.ServeMux)
	g.ICServers = make(map[string]*runtime.ServeMux)
	g.DataServerSeeds = make([]string, 0)
	g.ICServerSeeds = make([]string, 0)
	// create a new grpc server and register the inference gateway
	if config.Server.Test {
		go func() {

			gwMultiplexer, err := g.CreateServeMux("localhost:8090")
			if err != nil {
				slog.Error(fmt.Sprintf("register grpc restful api gateway handler error: %v", err))
			} else {
				g.DataServers["localhost:8090"] = gwMultiplexer
				g.ICServers["localhost:8090"] = gwMultiplexer
				grpcListener, _ := net.Listen("tcp", ":8090")
				grpcServer := grpc.NewServer(grpc.Creds(insecure.NewCredentials()))
				pb.RegisterGRPCInferenceServiceServer(grpcServer, new(Inference))
				if err := grpcServer.Serve(grpcListener); err != nil {
					slog.Error("grpc server start error: %v", err)
					panic(err)
				}
			}
		}()
	} else {
		for _, server := range config.DataServers {
			serverMux, err := g.CreateServeMux(server.Endpoint())
			if err == nil {
				g.DataServers[server.Endpoint()] = serverMux
				weight := server.Weight
				if weight < 0 {
					weight = 1
				}
				for weight > 0 {
					g.DataServerSeeds = append(g.DataServerSeeds, server.Endpoint())
					weight--
				}
			} else {
				slog.Error(fmt.Sprintf("create data server grpc restful api gateway handler error: %v", err))
			}
		}
		for _, server := range config.ICServers {
			serverMux, err := g.CreateServeMux(server.Endpoint())
			if err == nil {
				g.ICServers[server.Endpoint()] = serverMux
				weight := server.Weight
				if weight < 0 {
					weight = 1
				}
				for weight > 0 {
					g.ICServerSeeds = append(g.ICServerSeeds, server.Endpoint())
					weight--
				}
			} else {
				slog.Error(fmt.Sprintf("create ic server grpc restful api gateway handler error: %v", err))
			}
		}
	}
	slog.Info("ic servers.", "seeds", g.ICServerSeeds, "servers", g.ICServers)
	slog.Info("data server servers.", "seeds", g.DataServerSeeds, "servers", g.DataServers)
	if len(g.DataServers) <= 0 {
		panic(errors.New("no data servers"))
	}
	if len(g.ICServers) <= 0 {
		panic(errors.New("no ic servers"))
	}

	swaggerMux := swaggerui.NewServeMux(func(s string) ([]byte, error) {
		return []byte(pb.Swagger), nil
	}, "swagger.json")
	multiplexer.Handle("/swagger-ui/", swaggerMux)
	multiplexer.Handle("/swagger.json", swaggerMux)
	grpcMux := mux.NewChainMux()
	grpcMux.Use(auth.AuthenticationInterceptor)
	grpcMux.Use(ChargeInterceptor)
	grpcMux.HandleFunc("/comate/v1/grpc/dataserver/", func(writer http.ResponseWriter, request *http.Request) {
		g.PickDataServer().ServeHTTP(writer, request)
	})
	grpcMux.HandleFunc("/comate/v1/grpc/ic/", func(writer http.ResponseWriter, request *http.Request) {
		g.PickICServer().ServeHTTP(writer, request)
	})
	multiplexer.HandleFunc("/comate/v1/grpc/", func(writer http.ResponseWriter, request *http.Request) {
		grpcMux.ServeHTTP(writer, request)
	})
	return nil
}

func (g *GrpcServerAdapter) CreateServeMux(endpoint string) (*runtime.ServeMux, error) {
	gwMultiplexer := runtime.NewServeMux()
	err := pb.RegisterGRPCInferenceServiceHandlerFromEndpoint(context.Background(), gwMultiplexer,
		endpoint,
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())})
	return gwMultiplexer, err
}

func (g *GrpcServerAdapter) PickICServer() *runtime.ServeMux {
	pos := g.Rand.Intn(len(g.ICServerSeeds))
	slog.Info(fmt.Sprintf("pick ic server pos: %d", pos))
	return g.ICServers[g.ICServerSeeds[pos]]
}

func (g *GrpcServerAdapter) PickDataServer() *runtime.ServeMux {
	pos := g.Rand.Intn(len(g.DataServerSeeds))
	slog.Info(fmt.Sprintf("pick data server pos: %d", pos))
	return g.DataServers[g.DataServerSeeds[pos]]
}

func ChargeInterceptor(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cachedRequest, _ := mux.NewCachedRequest(r)
		response := &mux.CachedResponse{
			Request: cachedRequest,
			Writer:  w,
		}
		handler.ServeHTTP(response, r)
		go func() {
			conversation := NewRawConversation(response)
			_ = charge.ForConversation(conversation)
		}()
	})
}
