package telegram

import (
	"strings"
	"testing"

	"github.com/ujangdoubleday/gitel/internal/model"
)

func TestFormatMessage(t *testing.T) {
	event := &model.ExtractedPushEvent{
		Repository:  "ujangdoubleday/gitel",
		Pusher:      "ujangdoubleday",
		Branch:      "main",
		CommitCount: 2,
		Commits: []model.CommitInfo{
			{ID: "a1b2c3d", Author: "ujangdoubleday", Message: "feat: add webhook", Timestamp: "2024-01-01T10:00:00Z"},
		},
	}

	summary := "ringkasan:\n- menambahkan webhook receiver"
	result := FormatMessage(event, summary)

	if !strings.Contains(result, "<b>ujangdoubleday/gitel</b>") {
		t.Errorf("expected repository in bold, got: %s", result)
	}
	if !strings.Contains(result, "<code>main</code>") {
		t.Errorf("expected branch in code tag, got: %s", result)
	}
	if !strings.Contains(result, "menambahkan webhook receiver") {
		t.Errorf("expected summary text, got: %s", result)
	}
}

func TestFormatMessageEscapesHTML(t *testing.T) {
	event := &model.ExtractedPushEvent{
		Repository: "repo",
		Pusher:     "user",
		Branch:     "main",
	}

	summary := "<script>alert('xss')</script>"
	result := FormatMessage(event, summary)

	if strings.Contains(result, "<script>") {
		t.Errorf("HTML should be escaped, got: %s", result)
	}
	if !strings.Contains(result, "&lt;script&gt;") {
		t.Errorf("expected escaped script tag, got: %s", result)
	}
}

func TestFormatFallbackMessage(t *testing.T) {
	event := &model.ExtractedPushEvent{
		Repository:  "ujangdoubleday/gitel",
		Pusher:      "ujangdoubleday",
		Branch:      "dev",
		CommitCount: 1,
		Commits: []model.CommitInfo{
			{ID: "a1b2c3d", Author: "ujangdoubleday", Message: "fix: bug parser"},
		},
	}

	result := FormatFallbackMessage(event)

	if !strings.Contains(result, "<b>ujangdoubleday/gitel</b>") {
		t.Errorf("expected repository in bold, got: %s", result)
	}
	if !strings.Contains(result, "<code>a1b2c3d</code>") {
		t.Errorf("expected commit hash in code tag, got: %s", result)
	}
	if !strings.Contains(result, "fix: bug parser") {
		t.Errorf("expected commit message, got: %s", result)
	}
}
