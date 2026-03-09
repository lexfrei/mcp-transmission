package tools_test

import (
	"context"
	"testing"

	"github.com/cockroachdb/errors"
	"github.com/lexfrei/mcp-transmission/internal/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestQueueMoveTool_Definition(t *testing.T) {
	tool := tools.QueueMoveTool()

	if tool.Name != "transmission_queue_move" {
		t.Errorf("expected name transmission_queue_move, got %s", tool.Name)
	}
}

func TestQueueMoveHandler_Top(t *testing.T) {
	client := newMockClient()
	handler := tools.NewQueueMoveHandler(client)

	params := tools.QueueMoveParams{IDs: []int64{1}, Action: "top"}

	result, _, err := handler(context.Background(), &mcp.CallToolRequest{}, params)
	if err != nil {
		t.Fatalf("handler failed: %v", err)
	}

	if result != nil && result.IsError {
		t.Error("expected success")
	}
}

func TestQueueMoveHandler_InvalidAction(t *testing.T) {
	client := newMockClient()
	handler := tools.NewQueueMoveHandler(client)

	params := tools.QueueMoveParams{IDs: []int64{1}, Action: "invalid"}

	_, _, err := handler(context.Background(), &mcp.CallToolRequest{}, params)
	if !errors.Is(err, tools.ErrValidation) {
		t.Errorf("expected ErrValidation, got: %v", err)
	}
}

func TestQueueMoveHandler_MissingIDs(t *testing.T) {
	client := newMockClient()
	handler := tools.NewQueueMoveHandler(client)

	params := tools.QueueMoveParams{Action: "top"}

	_, _, err := handler(context.Background(), &mcp.CallToolRequest{}, params)
	if !errors.Is(err, tools.ErrValidation) {
		t.Errorf("expected ErrValidation, got: %v", err)
	}
}
