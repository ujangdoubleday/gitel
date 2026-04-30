package service

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/xoejang/gitel/internal/model"
	"github.com/xoejang/gitel/pkg/llm"
)

func TestFormatPushEvent(t *testing.T) {
	mockResponse := map[string]interface{}{
		"choices": []map[string]interface{}{
			{
				"message": map[string]interface{}{
					"role":    "assistant",
					"content": "ringkasan:\n- xoejang push ke branch main\n- menambahkan fitur webhook receiver\n- memperbaiki edge case parser",
				},
			},
		},
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/chat/completions" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("authorization") != "Bearer test-key" {
			t.Errorf("unexpected authorization header")
		}

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	client := llm.NewClient("test-key", "test-model", mockServer.URL, 10*time.Second)
	formatter := NewFormatter(client)

	event := &model.ExtractedPushEvent{
		Repository:  "xoejang/gitel",
		Pusher:      "xoejang",
		Branch:      "main",
		CommitCount: 2,
		Commits: []model.CommitInfo{
			{ID: "a1b2c3d", Author: "xoejang", Message: "feat: add webhook receiver", Timestamp: "2024-01-01T10:00:00Z", URL: "http://example.com/1"},
			{ID: "d4e5f6g", Author: "johndoe", Message: "fix: handle edge case", Timestamp: "2024-01-01T11:00:00Z", URL: "http://example.com/2"},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	summary, err := formatter.FormatPushEvent(ctx, event)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if summary == "" {
		t.Fatal("expected non-empty summary")
	}

	t.Logf("LLM summary: %s", summary)
}
