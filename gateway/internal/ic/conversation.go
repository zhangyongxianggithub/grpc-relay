package ic

import (
	"encoding/json"
	"strings"

	"bestzyx.com/grpc-relay/gateway/internal/charge"
	"bestzyx.com/grpc-relay/gateway/internal/mux"
)

type MsiConversation struct {
	Response    *mux.CachedResponse
	RequestJSON map[string]any
}

func (c *MsiConversation) UseModel() string {
	if model, ok := c.RequestJSON["model_id"].(string); ok {
		return model
	}
	return ""
}

func (c *MsiConversation) RequestId() string {
	if requestId, ok := c.RequestJSON["request_id"].(string); ok {
		return requestId
	}
	return ""
}

func (c *MsiConversation) Stream() bool {
	return true
}

func (c *MsiConversation) Ask() []string {
	if input, ok := c.RequestJSON["input"].(string); ok {
		return []string{input}
	}
	return []string{""}
}

func (c *MsiConversation) Reply() []string {
	reply := ""
	resp := string(c.Response.Response)
	events := strings.Split(resp, "\n\n")
	for _, str := range events {
		messages := strings.Split(str, "\n")
		for _, message := range messages {
			data := ""
			m := strings.TrimSpace(message)
			if d, found := strings.CutPrefix(m, "data:"); found {
				data += strings.TrimSpace(d)
			}
			if data != "" {
				msg := make(map[string]any)
				_ = json.Unmarshal([]byte(data), &msg)
				if result, ok := msg["output"].(string); ok {
					reply += result
				}
			}
		}
	}
	return []string{reply}
}

type RawConversation struct {
	Response    *mux.CachedResponse
	RequestJSON map[string]any
}

func (c *RawConversation) UseModel() string {
	if model, ok := c.RequestJSON["model_id"].(string); ok {
		return model
	}
	return ""
}

func (c *RawConversation) RequestId() string {
	if requestId, ok := c.RequestJSON["request_id"].(string); ok {
		return requestId
	}
	return ""
}

func (c *RawConversation) Stream() bool {
	return true
}

func (c *RawConversation) Ask() []string {
	return []string{c.Response.Request.String()}
}

func (c *RawConversation) Reply() []string {
	reply := make([]string, 0)
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
	return reply
}

func NewRawConversation(response *mux.CachedResponse) charge.Conversation {
	m, err := response.Request.JSON()
	if err != nil {
		return nil
	}
	return &RawConversation{Response: response, RequestJSON: m}
}

func NewMsiConversation(response *mux.CachedResponse) charge.Conversation {
	m, err := response.Request.JSON()
	if err != nil {
		return nil
	}
	return &MsiConversation{Response: response, RequestJSON: m}
}
