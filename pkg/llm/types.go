package llm

// message represents a single message in a chat completion request.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// chatRequest is the request body for OpenAI-compatible chat completions.
type chatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

// chatResponse is the response body from OpenAI-compatible chat completions.
type chatResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}
