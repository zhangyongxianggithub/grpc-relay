package charge

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"golang.org/x/xerrors"
)

type QaType string

const QUERY QaType = "query"
const ANSWER QaType = "answer"

type Record struct {
	QaType    QaType `json:"qaType"`
	ModelId   string `json:"modelId"`
	RequestId string `json:"requestId"`
	Content   string `json:"content"`
}

func ForTokens(record *Record) error {
	chargeContent, err := json.Marshal(record)
	if err != nil {
		slog.Error("marshal charge record failed",
			"qaType", record.QaType, "modelId", record.ModelId, "requestId", record.RequestId, "content",
			record.Content,
			slog.Any("error", xerrors.New(err.Error())))
		return nil
	}
	resp, err := service.client.Post(service.ChargeServer.GetUriPrefix()+"/record",
		"application/json",
		bytes.NewBuffer(chargeContent))
	if err != nil {
		slog.Error("request charge failed",
			"qaType", record.QaType, "modelId", record.ModelId, "requestId", record.RequestId, "content",
			record.Content,
			slog.Any("error", xerrors.New(err.Error())))
		return nil
	}
	response, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		slog.Error("request charge failed",
			"qaType", record.QaType, "modelId", record.ModelId, "requestId", record.RequestId, "content",
			record.Content, "status code", resp.StatusCode, "reason", string(response))
		return fmt.Errorf(string(response))
	} else {
		slog.Info("charge for conversation successfully.", "qaType", record.QaType, "modelId", record.ModelId,
			"requestId", record.RequestId, "response", string(response))
	}
	return nil
}

type Conversation interface {
	UseModel() string

	RequestId() string

	Stream() bool

	Ask() []string

	Reply() []string
}

func ForConversation(conversation Conversation) error {
	var err error
	for _, ask := range conversation.Ask() {
		err = ForTokens(&Record{
			QaType:    QUERY,
			ModelId:   conversation.UseModel(),
			RequestId: conversation.RequestId(),
			Content:   ask,
		})
	}
	for _, reply := range conversation.Reply() {
		err = ForTokens(&Record{
			QaType:    ANSWER,
			ModelId:   conversation.UseModel(),
			RequestId: conversation.RequestId(),
			Content:   reply,
		})
	}

	return err
}
