package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/ujangdoubleday/gitel/internal/model"
	"github.com/ujangdoubleday/gitel/pkg/llm"
)

// formatter builds prompts and calls the LLM to generate human-readable summaries.
type Formatter struct {
	client *llm.Client
}

// newFormatter creates a new formatter instance.
func NewFormatter(client *llm.Client) *Formatter {
	return &Formatter{client: client}
}

// formatPushEvent sends the extracted push event to the LLM and returns the formatted summary.
func (f *Formatter) FormatPushEvent(ctx context.Context, event *model.ExtractedPushEvent) (string, error) {
	prompt := buildPrompt(event)

	messages := []llm.Message{
		{
			Role:    "system",
			Content: "kamu adalah asisten yang merangkum commit git menjadi laporan yang mudah dibaca dalam bahasa indonesia. gunakan bahasa yang profesional namun santai.",
		},
		{
			Role:    "user",
			Content: prompt,
		},
	}

	return f.client.ChatCompletion(ctx, messages)
}

// buildPrompt constructs the prompt from extracted push event data.
func buildPrompt(event *model.ExtractedPushEvent) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("repository: %s\n", event.Repository))
	sb.WriteString(fmt.Sprintf("branch: %s\n", event.Branch))
	sb.WriteString(fmt.Sprintf("pusher: %s\n", event.Pusher))
	sb.WriteString(fmt.Sprintf("total commits: %d\n\n", event.CommitCount))

	sb.WriteString("daftar commit:\n")
	for i, c := range event.Commits {
		sb.WriteString(fmt.Sprintf("%d. [%s] %s - %s (%s)\n", i+1, c.ID[:7], c.Author, c.Message, c.Timestamp))
	}

	sb.WriteString("\nbuatkan ringkasan dalam bahasa indonesia yang mudah dibaca. sebutkan siapa yang push, branch apa, dan apa saja perubahan utama. format gunakan bullet point.")

	return sb.String()
}
