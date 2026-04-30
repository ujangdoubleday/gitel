package model

// pushEvent represents a GitHub push event webhook payload.
type PushEvent struct {
	Ref        string     `json:"ref"`
	Before     string     `json:"before"`
	After      string     `json:"after"`
	Repository Repository `json:"repository"`
	Pusher     Pusher     `json:"pusher"`
	Commits    []Commit   `json:"commits"`
	Compare    string     `json:"compare"`
}

// repository represents the repository information in a webhook payload.
type Repository struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	URL      string `json:"url"`
	Private  bool   `json:"private"`
}

// pusher represents the user who pushed the commits.
type Pusher struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// commit represents a single commit in the push event.
type Commit struct {
	ID        string   `json:"id"`
	Message   string   `json:"message"`
	Timestamp string   `json:"timestamp"`
	URL       string   `json:"url"`
	Author    Author   `json:"author"`
	Added     []string `json:"added"`
	Removed   []string `json:"removed"`
	Modified  []string `json:"modified"`
}

// author represents the author of a commit.
type Author struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
