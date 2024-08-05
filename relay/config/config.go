package config

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Initializer interface {
	Init(*ServerConfig) error
	Name() string
}

var initializers = make([]Initializer, 0)

var Server *ServerConfig

func AddInitializer(initializer Initializer) {
	initializers = append(initializers, initializer)
	slog.Info("register initializer.", "initializer", initializer.Name())
}
func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath("./config") // 还可以在工作目录中查找配置
	viper.AddConfigPath("./")
	viper.AddConfigPath("/etc/go/ic-relay/conf")
	err := viper.ReadInConfig() // 查找并读取配置文件
	if err != nil {
		if err, ok := err.(viper.ConfigFileNotFoundError); ok {
			slog.Error("config.toml file not found", slog.Any("err", err))
			panic(fmt.Errorf("config.toml file not found"))
		} else {
			slog.Error("read config.toml file failed", slog.Any("err", err))
			panic(err)
		}
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Printf("Config file %s changed, reload it", e.Name)
		InitializeConfig()
	})
}

func InitializeConfig() {
	if Server == nil {
		Server = new(ServerConfig)
	}
	// 反序列化配置文件
	err := viper.Unmarshal(Server)
	if Server.Gateway.Schema == "" {
		Server.Gateway.Schema = "http"
	}
	if err != nil {
		panic(fmt.Errorf("fail to load relay config file config.toml, err:%s \n", err))
	} else {
		config, _ := json.Marshal(Server)
		slog.Info("load relay config successfully.", "content", string(config))
	}
	for _, initializer := range initializers {
		err = initializer.Init(Server)
		if err != nil {
			panic(fmt.Errorf("%s initializer initialize failed", initializer.Name()))
		} else {
			slog.Info("initializer initialize successfully", "name", initializer.Name())
		}
	}
}

type ServerConfig struct {
	AppName string      `json:"appName"`
	Server  *GrpcServer `json:"server"`
	Gateway *Gateway    `json:"gateway"`
}

type GrpcServer struct {
	ICListen   string `json:"icListen"`
	DataListen string `json:"dataListen"`
	EBListen   string `json:"ebListen"`
}

type Gateway struct {
	Server string `json:"server"`
	Port   int    `json:"port"`
	Token  string `json:"token"`
	Schema string `json:"schema"`
}

func (g *Gateway) GetUriPrefix() string {
	return fmt.Sprintf("%s://%s", g.Schema, g.GetEndpoint())
}

func (g *Gateway) GetEndpoint() string {
	return fmt.Sprintf("%s:%d", g.Server, g.Port)
}
