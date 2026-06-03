package tools_test

import (
	"context"
	"strings"
	"testing"

	"github.com/cockroachdb/errors"
	"github.com/lexfrei/mcp-transmission/internal/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestTorrentRemoveTool_Definition(t *testing.T) {
	tool := tools.TorrentRemoveTool()

	if tool.Name != "transmission_torrent_remove" {
		t.Errorf("expected name transmission_torrent_remove, got %s", tool.Name)
	}
}

func TestTorrentRemoveHandler_Success(t *testing.T) {
	client := newMockClient()
	handler := tools.NewTorrentRemoveHandler(client)

	params := tools.TorrentRemoveParams{IDs: []int64{1, 2}}

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

func TestTorrentRemoveHandler_MissingIDs(t *testing.T) {
	client := newMockClient()
	handler := tools.NewTorrentRemoveHandler(client)

	params := tools.TorrentRemoveParams{}

	_, _, err := handler(context.Background(), &mcp.CallToolRequest{}, params)
	if !errors.Is(err, tools.ErrValidation) {
		t.Errorf("expected ErrValidation, got: %v", err)
	}
}

func TestTorrentRemoveHandler_DeleteLocalData_ElicitAccept(t *testing.T) {
	client := newMockClient()
	elicitor := &mockElicitor{
		result: &mcp.ElicitResult{Action: testActionAccept},
	}
	handler := tools.NewTorrentRemoveHandlerWithElicitor(client, elicitor)

	params := tools.TorrentRemoveParams{
		IDs:             []int64{1},
		DeleteLocalData: true,
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

	if !elicitor.called {
		t.Error("expected elicitation to be called")
	}
}

func TestTorrentRemoveHandler_DeleteLocalData_ElicitDecline(t *testing.T) {
	client := newMockClient()
	elicitor := &mockElicitor{
		result: &mcp.ElicitResult{Action: "decline"},
	}
	handler := tools.NewTorrentRemoveHandlerWithElicitor(client, elicitor)

	params := tools.TorrentRemoveParams{
		IDs:             []int64{1},
		DeleteLocalData: true,
	}

	result, output, err := handler(context.Background(), &mcp.CallToolRequest{}, params)
	if err != nil {
		t.Fatalf("handler failed: %v", err)
	}

	if result != nil && result.IsError {
		t.Error("expected no error flag on decline")
	}

	if !strings.Contains(strings.ToLower(output.Message), "cancelled") {
		t.Errorf("expected message to contain 'cancelled', got: %s", output.Message)
	}
}

func TestTorrentRemoveHandler_DeleteLocalData_ElicitCancel(t *testing.T) {
	client := newMockClient()
	elicitor := &mockElicitor{
		result: &mcp.ElicitResult{Action: "cancel"},
	}
	handler := tools.NewTorrentRemoveHandlerWithElicitor(client, elicitor)

	params := tools.TorrentRemoveParams{
		IDs:             []int64{1},
		DeleteLocalData: true,
	}

	result, output, err := handler(context.Background(), &mcp.CallToolRequest{}, params)
	if err != nil {
		t.Fatalf("handler failed: %v", err)
	}

	if result != nil && result.IsError {
		t.Error("expected no error flag on cancel")
	}

	if !strings.Contains(strings.ToLower(output.Message), "cancelled") {
		t.Errorf("expected message to contain 'cancelled', got: %s", output.Message)
	}
}

func TestTorrentRemoveHandler_DeleteLocalData_ElicitErrorFallbackConfirm(t *testing.T) {
	client := newMockClient()
	elicitor := &mockElicitor{err: errMock}
	handler := tools.NewTorrentRemoveHandlerWithElicitor(client, elicitor)

	// Elicitation fails but confirmDelete=true acts as fallback.
	params := tools.TorrentRemoveParams{
		IDs:             []int64{1},
		DeleteLocalData: true,
		ConfirmDelete:   true,
	}

	result, output, err := handler(context.Background(), &mcp.CallToolRequest{}, params)
	if err != nil {
		t.Fatalf("handler failed: %v", err)
	}

	if result != nil && result.IsError {
		t.Error("expected success with confirmDelete fallback")
	}

	if output.Message == "" {
		t.Error("expected non-empty message")
	}
}

func TestTorrentRemoveHandler_DeleteLocalData_ElicitErrorNoFallback(t *testing.T) {
	client := newMockClient()
	elicitor := &mockElicitor{err: errMock}
	handler := tools.NewTorrentRemoveHandlerWithElicitor(client, elicitor)

	// Elicitation fails and confirmDelete is false — should error.
	params := tools.TorrentRemoveParams{
		IDs:             []int64{1},
		DeleteLocalData: true,
	}

	_, _, err := handler(context.Background(), &mcp.CallToolRequest{}, params)
	if !errors.Is(err, tools.ErrValidation) {
		t.Errorf("expected ErrValidation (confirm required), got: %v", err)
	}
}

func TestTorrentRemoveHandler_DeleteLocalData_ElicitNilResult(t *testing.T) {
	client := newMockClient()
	elicitor := &mockElicitor{result: nil}
	handler := tools.NewTorrentRemoveHandlerWithElicitor(client, elicitor)

	// Nil result without confirmDelete fallback — should error.
	params := tools.TorrentRemoveParams{
		IDs:             []int64{1},
		DeleteLocalData: true,
	}

	_, _, err := handler(context.Background(), &mcp.CallToolRequest{}, params)
	if !errors.Is(err, tools.ErrValidation) {
		t.Errorf("expected ErrValidation for nil result without fallback, got: %v", err)
	}
}

func TestTorrentRemoveHandler_NoDeleteLocalData_SkipsElicitation(t *testing.T) {
	client := newMockClient()
	elicitor := &mockElicitor{
		result: &mcp.ElicitResult{Action: testActionAccept},
	}
	handler := tools.NewTorrentRemoveHandlerWithElicitor(client, elicitor)

	params := tools.TorrentRemoveParams{IDs: []int64{1}}

	_, _, err := handler(context.Background(), &mcp.CallToolRequest{}, params)
	if err != nil {
		t.Fatalf("handler failed: %v", err)
	}

	if elicitor.called {
		t.Error("elicitation should not be called when deleteLocalData is false")
	}
}

func TestTorrentRemoveHandler_DeleteLocalData_NilSessionWithConfirm(t *testing.T) {
	client := newMockClient()
	handler := tools.NewTorrentRemoveHandler(client)

	// No session (no elicitation), but confirmDelete=true as fallback.
	params := tools.TorrentRemoveParams{
		IDs:             []int64{1},
		DeleteLocalData: true,
		ConfirmDelete:   true,
	}

	result, output, err := handler(context.Background(), &mcp.CallToolRequest{}, params)
	if err != nil {
		t.Fatalf("handler failed: %v", err)
	}

	if result != nil && result.IsError {
		t.Error("expected success with confirmDelete fallback")
	}

	if output.Message == "" {
		t.Error("expected non-empty message")
	}
}

func TestTorrentRemoveHandler_DeleteLocalData_NilSessionNoConfirm(t *testing.T) {
	client := newMockClient()
	handler := tools.NewTorrentRemoveHandler(client)

	// No session, no confirmDelete — should error.
	params := tools.TorrentRemoveParams{
		IDs:             []int64{1},
		DeleteLocalData: true,
	}

	_, _, err := handler(context.Background(), &mcp.CallToolRequest{}, params)
	if !errors.Is(err, tools.ErrValidation) {
		t.Errorf("expected ErrValidation, got: %v", err)
	}
}

func TestTorrentRemoveHandler_DeleteLocalData_NilRequest(t *testing.T) {
	client := newMockClient()
	handler := tools.NewTorrentRemoveHandler(client)

	// Nil request, no confirmDelete — should error.
	params := tools.TorrentRemoveParams{
		IDs:             []int64{1},
		DeleteLocalData: true,
	}

	_, _, err := handler(context.Background(), nil, params)
	if !errors.Is(err, tools.ErrValidation) {
		t.Errorf("expected ErrValidation when request is nil, got: %v", err)
	}
}

func TestTorrentRemoveHandler_DeleteWithAccept_TransmissionError(t *testing.T) {
	client := newMockClient()
	client.err = errMock
	elicitor := &mockElicitor{result: &mcp.ElicitResult{Action: testActionAccept}}
	handler := tools.NewTorrentRemoveHandlerWithElicitor(client, elicitor)

	params := tools.TorrentRemoveParams{
		IDs:             []int64{1},
		DeleteLocalData: true,
	}

	_, _, err := handler(context.Background(), &mcp.CallToolRequest{}, params)
	if !errors.Is(err, tools.ErrTransmission) {
		t.Errorf("expected ErrTransmission, got: %v", err)
	}
}

func TestTorrentRemoveHandler_Error(t *testing.T) {
	client := newMockClient()
	client.err = errMock
	handler := tools.NewTorrentRemoveHandler(client)

	// deleteLocalData=false, so no elicitation needed — tests Transmission error path.
	params := tools.TorrentRemoveParams{IDs: []int64{1}}

	_, _, err := handler(context.Background(), &mcp.CallToolRequest{}, params)
	if !errors.Is(err, tools.ErrTransmission) {
		t.Errorf("expected ErrTransmission, got: %v", err)
	}
}

// mockElicitor implements tools.Elicitor for testing.
type mockElicitor struct {
	result *mcp.ElicitResult
	err    error
	called bool
}

func (m *mockElicitor) Elicit(_ context.Context, _ *mcp.ElicitParams) (*mcp.ElicitResult, error) {
	m.called = true

	if m.err != nil {
		return nil, m.err
	}

	return m.result, nil
}
