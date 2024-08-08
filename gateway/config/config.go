package config

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"github.com/zhangyongxianggithub/grpc-relay/gateway/internal/mux"
)

type Initializer interface {
	Init(*ServerConfig, *mux.ChainMux) error
	Name() string
}

var initializers = make([]Initializer, 0)

var Server *ServerConfig

func AddInitializer(initializer Initializer) {
	initializers = append(initializers, initializer)
}
func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath("./config") // 还可以在工作目录中查找配置
	viper.AddConfigPath("./")
	viper.AddConfigPath("/etc/go/ic-gateway/conf")
	viper.SetEnvKeyReplacer(strings.NewReplacer("..", "_"))
	err := viper.ReadInConfig() // 查找并读取配置文件
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			slog.Error("config.toml file not found", slog.Any("err", err))
			panic(fmt.Errorf("config.toml file not found"))
		} else {
			slog.Error("read config.toml file failed", slog.Any("err", err))
			panic(fmt.Errorf("read config.toml file failed, reason: %v\n", err))
		}
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		slog.Warn("config file changed, app should be restarted!")
	})
}

func InitializeConfig(multiplexer *mux.ChainMux) {
	if Server == nil {
		Server = new(ServerConfig)
	}
	err := viper.Unmarshal(Server)

	if err != nil {
		panic(fmt.Errorf("fail to load %s config file config.toml, err:%s \n", Server.AppName, err))
	} else {
		config, _ := json.Marshal(Server)
		slog.Info("load gateway config successfully.", "content", string(config))
	}
	if len(Server.ICServers) == 0 {
		panic("ic server config is empty")
	}
	if len(Server.DataServers) == 0 {
		panic("data server config is empty")
	}
	for _, initializer := range initializers {
		err = initializer.Init(Server, multiplexer)
		if err != nil {
			panic(fmt.Errorf("%s initializer initialize failed", initializer.Name()))
		} else {
			slog.Info("initializer initialize successfully", slog.String("name", initializer.Name()))
		}
	}
}

type ServerConfig struct {
	AppName       string            `json:"appName"`
	Server        *GatewayServer    `json:"server"`
	ICServers     []*UpstreamServer `json:"icServers"`
	DataServers   []*UpstreamServer `json:"dataServers"`
	ChargeService *UpstreamServer   `json:"chargeService"`
	EBMapping     map[string]string `json:"ebMapping"`
	Tokens        []string          `json:"tokens"`
}

type GatewayServer struct {
	Listen string `json:"listen"`
	Test   bool   `json:"test"`
}

type UpstreamServer struct {
	Name   string `json:"name"`
	Server string `json:"server"`
	Port   int    `json:"port"`
	Weight int    `json:"weight"`
}

func (s *UpstreamServer) Endpoint() string {
	return fmt.Sprintf("%s:%d", s.Server, s.Port)
}

func (s *UpstreamServer) GetUriPrefix() string {
	return fmt.Sprintf("http://%s", s.Endpoint())
}
