package service

import (
	"strings"

	"github.com/ujangdoubleday/gitel/internal/model"
)

// extractor handles parsing and filtering of GitHub webhook payloads.
type Extractor struct{}

// newExtractor creates a new extractor instance.
func NewExtractor() *Extractor {
	return &Extractor{}
}

// extract pulls out relevant information from a push event and filters commits.
func (e *Extractor) Extract(event *model.PushEvent) *model.ExtractedPushEvent {
	branch := strings.TrimPrefix(event.Ref, "refs/heads/")

	var commits []model.CommitInfo
	for _, c := range event.Commits {
		if e.shouldSkip(c) {
			continue
		}
		commits = append(commits, model.CommitInfo{
			ID:        c.ID,
			Author:    c.Author.Name,
			Message:   strings.TrimSpace(c.Message),
			Timestamp: c.Timestamp,
			URL:       c.URL,
		})
	}

	return &model.ExtractedPushEvent{
		Repository:  event.Repository.FullName,
		Pusher:      event.Pusher.Name,
		Branch:      branch,
		CommitCount: len(commits),
		Commits:     commits,
		CompareURL:  event.Compare,
	}
}

// shouldSkip determines if a commit should be excluded from the report.
func (e *Extractor) shouldSkip(c model.Commit) bool {
	// skip empty commits.
	if strings.TrimSpace(c.Message) == "" {
		return true
	}

	// skip merge commits.
	lowerMsg := strings.ToLower(c.Message)
	if strings.HasPrefix(lowerMsg, "merge ") {
		return true
	}

	return false
}
