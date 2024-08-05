package charge

import (
	"log/slog"
	"net/http"

	"bestzyx.com/grpc-relay/gateway/config"
	"bestzyx.com/grpc-relay/gateway/internal/mux"
)

type Service struct {
	ChargeServer *config.UpstreamServer
	client       *http.Client
}

var service = new(Service)

// init init 初始化函数，将 service 添加到 config 的 Initializer 中
func init() {
	config.AddInitializer(service)
	slog.Info("append charge-service to initializers")
}

// Init Init 初始化ChargeService服务，并返回错误信息（如果有）
func (c *Service) Init(serverConfig *config.ServerConfig, multiplexer *mux.ChainMux) error {
	c.client = &http.Client{
		Timeout: 0,
	}
	c.ChargeServer = serverConfig.ChargeService
	return nil
}

// Name func (c Service) Name() string
// 获取服务名称，实现ChargeService接口的Name方法
func (c *Service) Name() string {
	return "charge-service"
}
