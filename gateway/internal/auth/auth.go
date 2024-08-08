package auth

import (
	"net/http"

	"github.com/zhangyongxianggithub/grpc-relay/gateway/config"
	"github.com/zhangyongxianggithub/grpc-relay/gateway/internal/mux"
)

type Authentication struct {
	Tokens map[string]any
}

func (a *Authentication) Name() string {
	return "authentication"
}

var authentication *Authentication = new(Authentication)

func init() {
	config.AddInitializer(authentication)
}
func (a *Authentication) Init(config *config.ServerConfig, multiplexer *mux.ChainMux) error {
	a.Tokens = make(map[string]any)
	for _, token := range config.Tokens {
		a.Tokens[token] = true
	}
	return nil
}

func AuthenticationInterceptor(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		token := query.Get("token")
		if _, ok := authentication.Tokens[token]; !ok {
			w.WriteHeader(http.StatusUnauthorized)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"error":"Unauthorized"}`))
			return
		}
		handler.ServeHTTP(w, r)

	})
}
