package eb

import (
	"log/slog"
	"net/http"
	"strings"

	"bestzyx.com/grpc-relay/gateway/config"
	"bestzyx.com/grpc-relay/gateway/internal/charge"
	"bestzyx.com/grpc-relay/gateway/internal/mux"
	"github.com/google/uuid"
)

type ServiceProxy struct {
	Mappings map[string]string
}

func init() {
	config.AddInitializer(serviceProxy)
	slog.Info("append eb-service-proxy to initializers")
}

func (proxy *ServiceProxy) Init(serverConfig *config.ServerConfig, multiplexer *mux.ChainMux) error {
	proxy.Mappings = make(map[string]string)
	for url, model := range serverConfig.EBMapping {
		proxy.Mappings[strings.ReplaceAll(url, "!", ".")] = strings.ReplaceAll(model, "!", ".")
	}
	forwarder := NewForwarder()
	ebMux := mux.NewChainMux()
	ebMux.Use(RequestIdInterceptor)
	ebMux.Use(ChargeInterceptor)
	ebMux.HandleFunc("/", forwarder.Forward)
	multiplexer.HandleFunc("/rpc/2.0/ai_custom/", func(writer http.ResponseWriter, request *http.Request) {
		ebMux.ServeHTTP(writer, request)
	})
	multiplexer.HandleFunc("/oauth/2.0/token", forwarder.Forward)
	return nil
}

func RequestIdInterceptor(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Request-Id") == "" {
			requestId := uuid.New().String()
			slog.Info("request don't have a request id, generate a new one", "Request-Id", requestId)
			r.Header.Set("Request-Id", requestId)
		}
		handler.ServeHTTP(w, r)
	})
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
			if response.StatusCode == http.StatusOK {
				conversation := NewRawConversation(response)
				_ = charge.ForConversation(conversation)
			}
		}()
	})
}

func (proxy *ServiceProxy) Name() string {
	return "eb-service-proxy"
}

func (proxy *ServiceProxy) GetModel(uri string) string {
	if model, ok := proxy.Mappings[uri]; ok {
		return model
	} else {
		return "ERNIE-4.0-8K"
	}
}

var serviceProxy = new(ServiceProxy)
