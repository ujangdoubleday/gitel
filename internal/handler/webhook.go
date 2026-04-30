package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/ujangdoubleday/gitel/internal/model"
	"github.com/ujangdoubleday/gitel/internal/service"
)

const signatureHeader = "X-Hub-Signature-256"
const signaturePrefix = "sha256="

// webhookHandler handles GitHub webhook events.
type WebhookHandler struct {
	secret    string
	extractor *service.Extractor
	processor *service.Processor
}

// newWebhookHandler creates a new webhookHandler.
func NewWebhookHandler(secret string, extractor *service.Extractor, processor *service.Processor) *WebhookHandler {
	return &WebhookHandler{
		secret:    secret,
		extractor: extractor,
		processor: processor,
	}
}

// handleGitHubWebhook receives and processes GitHub push events asynchronously.
func (h *WebhookHandler) HandleGitHubWebhook(w http.ResponseWriter, r *http.Request) {
	// limit request body to 1MB to prevent abuse.
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("failed to read webhook body", "error", err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	// verify HMAC-SHA256 signature.
	if err := h.verifySignature(r, body); err != nil {
		slog.Error("signature verification failed", "error", err)
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// parse the push event.
	var event model.PushEvent
	if err := json.Unmarshal(body, &event); err != nil {
		slog.Error("failed to parse webhook JSON", "error", err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	extracted := h.extractor.Extract(&event)

	slog.Info("webhook received",
		"repository", extracted.Repository,
		"branch", extracted.Branch,
		"pusher", extracted.Pusher,
		"commits_after_filter", extracted.CommitCount,
	)

	// process LLM and telegram asynchronously so github gets immediate response.
	h.processor.ProcessAsync(extracted)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

// verifySignature validates the X-Hub-Signature-256 header.
func (h *WebhookHandler) verifySignature(r *http.Request, body []byte) error {
	sig := r.Header.Get(signatureHeader)
	if sig == "" {
		return errors.New("missing signature header")
	}

	if len(sig) <= len(signaturePrefix) || sig[:len(signaturePrefix)] != signaturePrefix {
		return errors.New("invalid signature format")
	}

	mac := hmac.New(sha256.New, []byte(h.secret))
	mac.Write(body)
	expectedMAC := mac.Sum(nil)

	expectedSig, err := hex.DecodeString(sig[len(signaturePrefix):])
	if err != nil {
		return errors.New("failed to decode signature")
	}

	if !hmac.Equal(expectedMAC, expectedSig) {
		return errors.New("signature mismatch")
	}

	return nil
}
