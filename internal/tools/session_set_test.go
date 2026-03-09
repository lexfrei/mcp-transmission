package tools_test

import (
	"context"
	"testing"

	"github.com/cockroachdb/errors"
	"github.com/lexfrei/mcp-transmission/internal/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestSessionSetTool_Definition(t *testing.T) {
	tool := tools.SessionSetTool()

	if tool.Name != "transmission_session_set" {
		t.Errorf("expected name transmission_session_set, got %s", tool.Name)
	}
}

func TestSessionSetHandler_Success(t *testing.T) {
	client := newMockClient()
	handler := tools.NewSessionSetHandler(client)

	limit := int64(1024)
	params := tools.SessionSetParams{
		SpeedLimitDown: &limit,
	}

	result, output, err := handler(context.Background(), &mcp.CallToolRequest{}, params)
	if err != nil {
		t.Fatalf("handler failed: %v", err)
	}

	if result != nil && result.IsError {
		t.Error("expected success")
	}

	if output.Message == "" {
		t.Error("expected non-empty message")
	}
}

func TestSessionSetHandler_EmptyParams(t *testing.T) {
	client := newMockClient()
	handler := tools.NewSessionSetHandler(client)

	params := tools.SessionSetParams{}

	_, _, err := handler(context.Background(), &mcp.CallToolRequest{}, params)
	if !errors.Is(err, tools.ErrValidation) {
		t.Errorf("expected ErrValidation for empty params, got: %v", err)
	}
}

func TestSessionSetHandler_TransmissionError(t *testing.T) {
	client := newMockClient()
	client.err = errMock

	handler := tools.NewSessionSetHandler(client)

	limit := int64(1024)
	params := tools.SessionSetParams{
		SpeedLimitDown: &limit,
	}

	_, _, err := handler(context.Background(), &mcp.CallToolRequest{}, params)
	if !errors.Is(err, tools.ErrTransmission) {
		t.Errorf("expected ErrTransmission, got: %v", err)
	}
}
