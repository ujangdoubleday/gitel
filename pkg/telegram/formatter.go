package telegram

import (
	"fmt"
	"html"
	"strings"

	"github.com/xoejang/gitel/internal/model"
)

// formatMessage builds an HTML-formatted Telegram message from the LLM summary.
func FormatMessage(event *model.ExtractedPushEvent, summary string) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("<b>%s</b>\n", html.EscapeString(event.Repository)))
	sb.WriteString(fmt.Sprintf("branch: <code>%s</code> | pusher: %s\n\n",
		html.EscapeString(event.Branch), html.EscapeString(event.Pusher)))
	sb.WriteString(html.EscapeString(summary))

	return sb.String()
}

// formatFallbackMessage builds a simple HTML message without LLM summary.
func FormatFallbackMessage(event *model.ExtractedPushEvent) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("<b>%s</b>\n", html.EscapeString(event.Repository)))
	sb.WriteString(fmt.Sprintf("branch: <code>%s</code> | pusher: %s\n\n",
		html.EscapeString(event.Branch), html.EscapeString(event.Pusher)))

	sb.WriteString("commits:\n")
	for _, c := range event.Commits {
		sb.WriteString(fmt.Sprintf("- <code>%s</code> %s: %s\n",
			html.EscapeString(c.ID[:7]),
			html.EscapeString(c.Author),
			html.EscapeString(c.Message)))
	}

	return sb.String()
}
