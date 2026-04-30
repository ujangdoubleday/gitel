package service

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/ujangdoubleday/gitel/internal/model"
	"github.com/ujangdoubleday/gitel/pkg/telegram"
)

// processor handles the async pipeline from extracted event to telegram message.
type Processor struct {
	formatter *Formatter
	telegram  *telegram.Client
	wg        sync.WaitGroup
}

// newProcessor creates a new processor instance.
func NewProcessor(formatter *Formatter, telegram *telegram.Client) *Processor {
	return &Processor{
		formatter: formatter,
		telegram:  telegram,
	}
}

// processAsync starts the LLM formatting and telegram sending in a background goroutine.
func (p *Processor) ProcessAsync(event *model.ExtractedPushEvent) {
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		p.run(event)
	}()
}

// run executes the synchronous pipeline: LLM -> telegram.
func (p *Processor) run(event *model.ExtractedPushEvent) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	slog.Info("starting LLM formatting",
		"repository", event.Repository,
		"branch", event.Branch,
		"commit_count", event.CommitCount,
	)

	summary, err := p.formatter.FormatPushEvent(ctx, event)
	if err != nil {
		slog.Error("LLM formatting failed",
			"error", err,
			"repository", event.Repository,
		)
		msg := telegram.FormatFallbackMessage(event)
		if err := p.telegram.SendMessage(ctx, msg); err != nil {
			slog.Error("telegram fallback send failed",
				"error", err,
				"repository", event.Repository,
			)
		} else {
			slog.Info("telegram fallback sent",
				"repository", event.Repository,
			)
		}
		return
	}

	slog.Info("LLM summary generated",
		"repository", event.Repository,
		"summary_length", len(summary),
	)

	msg := telegram.FormatMessage(event, summary)
	if err := p.telegram.SendMessage(ctx, msg); err != nil {
		slog.Error("telegram send failed",
			"error", err,
			"repository", event.Repository,
		)
		return
	}

	slog.Info("telegram message sent",
		"repository", event.Repository,
		"branch", event.Branch,
	)
}

// wait blocks until all background jobs finish.
func (p *Processor) Wait() {
	p.wg.Wait()
}
