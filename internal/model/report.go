package model

// extractedPushEvent holds the cleaned and filtered data from a GitHub push event.
type ExtractedPushEvent struct {
	Repository  string
	Pusher      string
	Branch      string
	CommitCount int
	Commits     []CommitInfo
	CompareURL  string
}

// commitInfo holds simplified commit data for reporting.
type CommitInfo struct {
	ID        string
	Author    string
	Message   string
	Timestamp string
	URL       string
}
