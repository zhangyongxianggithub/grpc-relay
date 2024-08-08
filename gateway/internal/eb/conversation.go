package eb

import (
	"encoding/json"
	"log/slog"
	"strings"

	"github.com/zhangyongxianggithub/grpc-relay/gateway/internal/charge"
	"github.com/zhangyongxianggithub/grpc-relay/gateway/internal/mux"
)

type Conversation struct {
	Response    *mux.CachedResponse
	RequestJson map[string]any
}

func (c *Conversation) UseModel() string {
	return serviceProxy.GetModel(c.Response.Request.Path())
}

func (c *Conversation) RequestId() string {
	return c.Response.Request.Header("Request-Id")
}

func (c *Conversation) Stream() bool {
	m := c.RequestJson
	if stream, ok := m["stream"].(bool); ok {
		return stream
	}
	return false
}

func (c *Conversation) Ask() []string {
	msgs := make([]string, 0)
	m := c.RequestJson
	if messages, ok := m["messages"].([]any); ok {
		for _, message := range messages {
			if message, ok := message.(map[string]interface{}); ok {
				if content, ok := message["content"].(string); ok && content != "" {
					msgs = append(msgs, content)
				}
			}
		}
	}
	return []string{strings.Join(msgs, " ")}
}

func (c *Conversation) Reply() []string {
	if !c.Stream() {
		m := map[string]any{}
		_ = json.Unmarshal(c.Response.Response, &m)
		if result, ok := m["result"].(string); ok {
			return []string{result}
		}
	} else {
		reply := ""
		resp := string(c.Response.Response)
		events := strings.Split(resp, "\n\n")
		for _, str := range events {
			messages := strings.Split(str, "\n")
			data := ""
			for _, message := range messages {
				m := strings.TrimSpace(message)
				if d, found := strings.CutPrefix(m, "data:"); found {
					data += strings.TrimSpace(d)
				}
			}
			if data != "" {
				msg := make(map[string]any)
				_ = json.Unmarshal([]byte(data), &msg)
				if result, ok := msg["result"].(string); ok {
					reply += result
				}
			}
			// slog.Info(data)
		}
		return []string{reply}
	}
	return []string{""}
}

func NewConversation(response *mux.CachedResponse) charge.Conversation {
	m, err := response.Request.JSON()
	if err != nil {
		slog.Error("unmarshal request failed.", slog.Any("err", err))
		m = make(map[string]any)
	}
	return &Conversation{
		Response:    response,
		RequestJson: m,
	}
}

type RawConversation struct {
	Response    *mux.CachedResponse
	RequestJSON map[string]any
}

func (c *RawConversation) UseModel() string {
	return serviceProxy.GetModel(c.Response.Request.Path())
}

func (c *RawConversation) RequestId() string {
	return c.Response.Request.Header("Request-Id")
}

func (c *RawConversation) Stream() bool {
	m := c.RequestJSON
	if stream, ok := m["stream"].(bool); ok {
		return stream
	}
	return false
}

func (c *RawConversation) Ask() []string {
	return []string{c.Response.Request.String()}
}

func (c *RawConversation) Reply() []string {
	reply := make([]string, 0)
	if !c.Stream() {
		reply = append(reply, string(c.Response.Response))
	} else {
		resp := string(c.Response.Response)
		events := strings.Split(resp, "\n\n")
		for _, str := range events {
			messages := strings.Split(str, "\n")
			data := ""
			for _, message := range messages {
				m := strings.TrimSpace(message)
				if d, found := strings.CutPrefix(m, "data:"); found {
					data += strings.TrimSpace(d)
				}
			}
			if data != "" {
				reply = append(reply, data)
			}
		}
	}
	return reply
}

func NewRawConversation(response *mux.CachedResponse) charge.Conversation {
	m, err := response.Request.JSON()
	if err != nil {
		return nil
	}
	return &RawConversation{Response: response, RequestJSON: m}
}
